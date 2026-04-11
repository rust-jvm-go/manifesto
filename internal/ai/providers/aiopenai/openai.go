package aiopenai

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Abraxas-365/manifesto/internal/ai/embedding"
	"github.com/Abraxas-365/manifesto/internal/ai/llm"
	"github.com/Abraxas-365/manifesto/internal/ai/speech"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
	"github.com/openai/openai-go/v3/shared"
)

// OpenAIProvider implements the LLM interface for OpenAI
type OpenAIProvider struct {
	client openai.Client
	apiKey string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string, opts ...option.RequestOption) *OpenAIProvider {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	options := append([]option.RequestOption{option.WithAPIKey(apiKey)}, opts...)
	client := openai.NewClient(options...)

	return &OpenAIProvider{
		client: client,
		apiKey: apiKey,
	}
}

func defaultChatOptions() *llm.ChatOptions {
	options := llm.DefaultOptions()
	options.Model = "gpt-4.1"
	return options
}

// modelSupportsTemperature returns false for models that reject the temperature
// parameter on the Responses API: o1/o3/o4 reasoning models, codex, and computer-use.
func modelSupportsTemperature(model string) bool {
	noTempPrefixes := []string{"o1", "o3", "o4", "codex", "computer-use"}
	for _, prefix := range noTempPrefixes {
		if strings.HasPrefix(model, prefix) {
			return false
		}
	}
	return true
}

// ============================================================================
// Chat Implementation
// ============================================================================

// Chat implements the LLM interface
func (p *OpenAIProvider) Chat(ctx context.Context, messages []llm.Message, opts ...llm.Option) (llm.Response, error) {
	if p.apiKey == "" {
		return llm.Response{}, errorRegistry.New(ErrMissingAPIKey)
	}
	if len(messages) == 0 {
		return llm.Response{}, errorRegistry.New(ErrEmptyMessages)
	}

	options := defaultChatOptions()
	for _, opt := range opts {
		opt(options)
	}

	instructions, inputItems, err := convertMessagesToResponsesInput(messages)
	if err != nil {
		return llm.Response{}, WrapError(err, ErrInvalidMessage).WithDetail("error", "failed to convert messages")
	}

	params := responses.ResponseNewParams{
		Model:                options.Model,
		Store:                openai.Bool(true),
		PromptCacheRetention: responses.ResponseNewParamsPromptCacheRetention24h,
		Truncation:           responses.ResponseNewParamsTruncationAuto,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: inputItems,
		},
	}

	if instructions != "" {
		params.Instructions = openai.String(instructions)
	}
	if options.Temperature != 0 && modelSupportsTemperature(options.Model) {
		params.Temperature = openai.Float(float64(options.Temperature))
	}
	if options.TopP != 0 && modelSupportsTemperature(options.Model) {
		params.TopP = openai.Float(float64(options.TopP))
	}
	if options.MaxCompletionTokens > 0 {
		params.MaxOutputTokens = openai.Int(int64(options.MaxCompletionTokens))
	} else if options.MaxTokens > 0 {
		params.MaxOutputTokens = openai.Int(int64(options.MaxTokens))
	}
	if options.User != "" {
		params.User = openai.String(options.User)
	}
	if options.ReasoningEffort != "" {
		params.Reasoning = shared.ReasoningParam{
			Effort: convertToReasoningEffort(options.ReasoningEffort),
		}
	}
	if len(options.Tools) > 0 || len(options.Functions) > 0 {
		tools, err := convertToResponsesTools(options.Tools, options.Functions)
		if err != nil {
			return llm.Response{}, WrapError(err, ErrConversionFailed).WithDetail("error", "failed to convert tools")
		}
		if len(tools) > 0 {
			params.Tools = tools
		}
	}
	if options.ToolChoice != nil {
		params.ToolChoice = convertToResponsesToolChoice(options.ToolChoice)
	}
	if options.JSONMode {
		params.Text = responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{},
			},
		}
	} else if options.ResponseFormat != nil {
		textFormat, err := convertToResponsesFormatParam(options.ResponseFormat)
		if err != nil {
			return llm.Response{}, WrapError(err, ErrConversionFailed).WithDetail("error", "failed to convert response format")
		}
		params.Text = responses.ResponseTextConfigParam{Format: textFormat}
	}

	resp, err := p.client.Responses.New(ctx, params)
	if err != nil {
		return llm.Response{}, ParseOpenAIError(err).
			WithDetail("model", options.Model).
			WithDetail("num_messages", len(messages))
	}

	result, err := convertFromResponsesResponse(resp)
	if err != nil {
		return llm.Response{}, WrapError(err, ErrAPIResponse).WithDetail("error", "failed to parse response")
	}
	return result, nil
}

// ============================================================================
// Chat Stream Implementation
// ============================================================================

// ChatStream implements streaming for Responses API
func (p *OpenAIProvider) ChatStream(ctx context.Context, messages []llm.Message, opts ...llm.Option) (llm.Stream, error) {
	if p.apiKey == "" {
		return nil, errorRegistry.New(ErrMissingAPIKey)
	}
	if len(messages) == 0 {
		return nil, errorRegistry.New(ErrEmptyMessages)
	}

	options := defaultChatOptions()
	for _, opt := range opts {
		opt(options)
	}

	instructions, inputItems, err := convertMessagesToResponsesInput(messages)
	if err != nil {
		return nil, WrapError(err, ErrInvalidMessage).WithDetail("error", "failed to convert messages")
	}

	params := responses.ResponseNewParams{
		Model:                options.Model,
		Store:                openai.Bool(true),
		PromptCacheRetention: responses.ResponseNewParamsPromptCacheRetention24h,
		Truncation:           responses.ResponseNewParamsTruncationAuto,
		Input: responses.ResponseNewParamsInputUnion{
			OfInputItemList: inputItems,
		},
	}

	if instructions != "" {
		params.Instructions = openai.String(instructions)
	}
	if options.Temperature != 0 && modelSupportsTemperature(options.Model) {
		params.Temperature = openai.Float(float64(options.Temperature))
	}
	if options.TopP != 0 && modelSupportsTemperature(options.Model) {
		params.TopP = openai.Float(float64(options.TopP))
	}
	if options.MaxCompletionTokens > 0 {
		params.MaxOutputTokens = openai.Int(int64(options.MaxCompletionTokens))
	} else if options.MaxTokens > 0 {
		params.MaxOutputTokens = openai.Int(int64(options.MaxTokens))
	}
	if options.User != "" {
		params.User = openai.String(options.User)
	}
	if options.ReasoningEffort != "" {
		params.Reasoning = shared.ReasoningParam{
			Effort: convertToReasoningEffort(options.ReasoningEffort),
		}
	}
	if len(options.Tools) > 0 || len(options.Functions) > 0 {
		tools, err := convertToResponsesTools(options.Tools, options.Functions)
		if err != nil {
			return nil, WrapError(err, ErrConversionFailed)
		}
		if len(tools) > 0 {
			params.Tools = tools
		}
	}
	if options.ToolChoice != nil {
		params.ToolChoice = convertToResponsesToolChoice(options.ToolChoice)
	}
	if options.JSONMode {
		params.Text = responses.ResponseTextConfigParam{
			Format: responses.ResponseFormatTextConfigUnionParam{
				OfJSONObject: &shared.ResponseFormatJSONObjectParam{},
			},
		}
	} else if options.ResponseFormat != nil {
		textFormat, err := convertToResponsesFormatParam(options.ResponseFormat)
		if err != nil {
			return nil, WrapError(err, ErrConversionFailed).WithDetail("error", "failed to convert response format")
		}
		params.Text = responses.ResponseTextConfigParam{Format: textFormat}
	}

	sseStream := p.client.Responses.NewStreaming(ctx, params)
	return &openAIStream{stream: sseStream}, nil
}

// ============================================================================
// Embedding Implementation
// ============================================================================

func (p *OpenAIProvider) EmbedDocuments(ctx context.Context, documents []string, opts ...embedding.Option) ([]embedding.Embedding, error) {
	// Validate input
	if len(documents) == 0 {
		return nil, errorRegistry.New(ErrEmptyEmbeddingInput)
	}

	options := embedding.DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	params := openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfArrayOfStrings: documents,
		},
	}

	if options.Model != "" {
		params.Model = options.Model
	} else {
		params.Model = "text-embedding-3-small"
	}

	if options.Dimensions > 0 {
		params.Dimensions = openai.Int(int64(options.Dimensions))
	}

	if options.User != "" {
		params.User = openai.String(options.User)
	}

	resp, err := p.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, ParseOpenAIError(err).
			WithDetail("model", params.Model).
			WithDetail("num_documents", len(documents))
	}

	if len(resp.Data) == 0 {
		return nil, errorRegistry.New(ErrNoEmbeddingReturned).
			WithDetail("num_documents", len(documents))
	}

	embeddings := make([]embedding.Embedding, len(resp.Data))
	for i, data := range resp.Data {
		embeddings[i] = embedding.Embedding{
			Vector: convertToFloat32Slice(data.Embedding),
			Usage: embedding.Usage{
				PromptTokens: int(resp.Usage.PromptTokens),
				TotalTokens:  int(resp.Usage.TotalTokens),
			},
		}
	}

	return embeddings, nil
}

func (p *OpenAIProvider) EmbedQuery(ctx context.Context, text string, opts ...embedding.Option) (embedding.Embedding, error) {
	if text == "" {
		return embedding.Embedding{}, errorRegistry.New(ErrEmptyEmbeddingInput)
	}

	embeddings, err := p.EmbedDocuments(ctx, []string{text}, opts...)
	if err != nil {
		return embedding.Embedding{}, err
	}

	if len(embeddings) == 0 {
		return embedding.Embedding{}, errorRegistry.New(ErrNoEmbeddingReturned)
	}

	return embeddings[0], nil
}

// ============================================================================
// Speech Synthesis Implementation
// ============================================================================

func (p *OpenAIProvider) Synthesize(ctx context.Context, text string, opts ...speech.SynthesisOption) (speech.Audio, error) {
	if text == "" {
		return speech.Audio{}, errorRegistry.New(ErrEmptySpeechInput)
	}

	options := speech.SynthesisOptions{
		Model:       string(openai.SpeechModelTTS1),
		Voice:       "alloy",
		AudioFormat: speech.AudioFormatMP3,
		SpeechRate:  1.0,
	}

	for _, opt := range opts {
		opt(&options)
	}

	responseFormat := openai.AudioSpeechNewParamsResponseFormatMP3
	switch options.AudioFormat {
	case speech.AudioFormatMP3:
		responseFormat = openai.AudioSpeechNewParamsResponseFormatMP3
	case speech.AudioFormatPCM:
		responseFormat = openai.AudioSpeechNewParamsResponseFormatPCM
	case speech.AudioFormatOGG:
		responseFormat = openai.AudioSpeechNewParamsResponseFormatOpus
	case speech.AudioFormatWAV:
		responseFormat = openai.AudioSpeechNewParamsResponseFormatPCM
	}

	voice := openai.AudioSpeechNewParamsVoiceAlloy
	switch strings.ToLower(options.Voice) {
	case "alloy":
		voice = openai.AudioSpeechNewParamsVoiceAlloy
	case "echo":
		voice = openai.AudioSpeechNewParamsVoiceEcho
	case "fable":
	case "shimmer":
		voice = openai.AudioSpeechNewParamsVoiceShimmer
	default:
		voice = openai.AudioSpeechNewParamsVoiceAlloy
	}

	params := openai.AudioSpeechNewParams{
		Model:          options.Model,
		Input:          text,
		Voice:          voice,
		ResponseFormat: responseFormat,
	}

	if options.SpeechRate != 1.0 {
		params.Speed = param.NewOpt(float64(options.SpeechRate))
	}

	res, err := p.client.Audio.Speech.New(ctx, params)
	if err != nil {
		return speech.Audio{}, ParseOpenAIError(err).
			WithDetail("model", options.Model).
			WithDetail("voice", options.Voice)
	}

	sampleRate := 24000
	if options.SampleRate > 0 {
		sampleRate = options.SampleRate
	}

	return speech.Audio{
		Content:    res.Body,
		Format:     options.AudioFormat,
		SampleRate: sampleRate,
		Usage: speech.TTSUsage{
			InputCharacters: len(text),
		},
	}, nil
}

// ============================================================================
// Speech Transcription Implementation
// ============================================================================

func (p *OpenAIProvider) Transcribe(ctx context.Context, audio io.Reader, opts ...speech.TranscriptionOption) (speech.Transcript, error) {
	if audio == nil {
		return speech.Transcript{}, errorRegistry.New(ErrEmptySpeechInput).
			WithDetail("error", "audio reader cannot be nil")
	}

	options := speech.TranscriptionOptions{
		Model:      string(openai.AudioModelWhisper1),
		Language:   "",
		Timestamps: false,
	}

	for _, opt := range opts {
		opt(&options)
	}

	params := openai.AudioTranscriptionNewParams{
		Model: options.Model,
		File:  audio,
	}

	if options.Language != "" {
		params.Language = param.NewOpt(options.Language)
	}

	response, err := p.client.Audio.Transcriptions.New(ctx, params)
	if err != nil {
		return speech.Transcript{}, ParseOpenAIError(err).
			WithDetail("model", options.Model)
	}

	result := speech.Transcript{
		Text: response.Text,
	}

	return result, nil
}

// ============================================================================
// Stream Implementation
// ============================================================================

type openAIStream struct {
	stream interface {
		Next() bool
		Current() responses.ResponseStreamEventUnion
		Err() error
		Close() error
	}
	done bool
}

func (s *openAIStream) Next() (llm.Message, error) {
	if s.done {
		return llm.Message{}, io.EOF
	}
	for s.stream.Next() {
		event := s.stream.Current()
		switch event.Type {
		case "response.output_text.delta":
			return llm.Message{Role: llm.RoleAssistant, Content: event.Delta}, nil
		case "response.output_item.done":
			item := event.Item
			if item.Type == "function_call" {
				fc := item.AsFunctionCall()
				return llm.Message{
					Role: llm.RoleAssistant,
					ToolCalls: []llm.ToolCall{{
						ID:   fc.CallID,
						Type: "function",
						Function: llm.FunctionCall{Name: fc.Name, Arguments: fc.Arguments},
					}},
				}, nil
			}
		case "response.failed":
			s.done = true
			return llm.Message{}, errorRegistry.New(ErrAPIResponse).WithDetail("error", "response failed")
		}
	}
	s.done = true
	if err := s.stream.Err(); err != nil {
		return llm.Message{}, ParseOpenAIError(err)
	}
	return llm.Message{}, io.EOF
}

func (s *openAIStream) Close() error {
	return s.stream.Close()
}

// ============================================================================
// Helper Functions
// ============================================================================

func convertMessagesToResponsesInput(messages []llm.Message) (string, responses.ResponseInputParam, error) {
	var systemParts []string
	inputItems := make(responses.ResponseInputParam, 0, len(messages))

	for i, msg := range messages {
		switch msg.Role {
		case llm.RoleSystem:
			if msg.Content != "" {
				systemParts = append(systemParts, msg.Content)
			}
		case llm.RoleUser:
			if msg.IsMultimodal() {
				contentList, err := convertToResponsesContentParts(msg.MultiContent)
				if err != nil {
					return "", nil, WrapError(err, ErrInvalidMessage).WithDetail("message_index", i)
				}
				inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
					OfMessage: &responses.EasyInputMessageParam{
						Role: responses.EasyInputMessageRoleUser,
						Content: responses.EasyInputMessageContentUnionParam{
							OfInputItemContentList: contentList,
						},
					},
				})
			} else {
				inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
					OfMessage: &responses.EasyInputMessageParam{
						Role: responses.EasyInputMessageRoleUser,
						Content: responses.EasyInputMessageContentUnionParam{
							OfString: openai.String(msg.Content),
						},
					},
				})
			}
		case llm.RoleAssistant:
			if len(msg.ToolCalls) > 0 {
				if msg.Content != "" {
					inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
						OfMessage: &responses.EasyInputMessageParam{
							Role: responses.EasyInputMessageRoleAssistant,
							Content: responses.EasyInputMessageContentUnionParam{
								OfString: openai.String(msg.Content),
							},
						},
					})
				}
				for _, tc := range msg.ToolCalls {
					inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
						OfFunctionCall: &responses.ResponseFunctionToolCallParam{
							CallID:    tc.ID,
							Name:      tc.Function.Name,
							Arguments: tc.Function.Arguments,
						},
					})
				}
			} else {
				inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
					OfMessage: &responses.EasyInputMessageParam{
						Role: responses.EasyInputMessageRoleAssistant,
						Content: responses.EasyInputMessageContentUnionParam{
							OfString: openai.String(msg.Content),
						},
					},
				})
			}
		case llm.RoleTool:
			inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
				OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
					CallID: msg.ToolCallID,
					Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
						OfString: openai.String(msg.Content),
					},
				},
			})
		case llm.RoleFunction:
			inputItems = append(inputItems, responses.ResponseInputItemUnionParam{
				OfFunctionCallOutput: &responses.ResponseInputItemFunctionCallOutputParam{
					CallID: msg.Name,
					Output: responses.ResponseInputItemFunctionCallOutputOutputUnionParam{
						OfString: openai.String(msg.Content),
					},
				},
			})
		default:
			return "", nil, errorRegistry.New(ErrUnsupportedRole).WithDetail("role", msg.Role)
		}
	}

	return strings.Join(systemParts, "\n\n"), inputItems, nil
}

func convertToResponsesTools(tools []llm.Tool, functions []llm.Function) ([]responses.ToolUnionParam, error) {
	result := make([]responses.ToolUnionParam, 0, len(tools)+len(functions))
	for _, tool := range tools {
		if tool.Type != "function" {
			continue
		}
		paramsMap, err := toParamsMap(tool.Function.Parameters, tool.Function.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        tool.Function.Name,
				Description: openai.String(tool.Function.Description),
				Parameters:  paramsMap,
			},
		})
	}
	for _, fn := range functions {
		paramsMap, err := toParamsMap(fn.Parameters, fn.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, responses.ToolUnionParam{
			OfFunction: &responses.FunctionToolParam{
				Name:        fn.Name,
				Description: openai.String(fn.Description),
				Parameters:  paramsMap,
			},
		})
	}
	return result, nil
}

func toParamsMap(parameters any, name string) (map[string]any, error) {
	paramsJSON, err := json.Marshal(parameters)
	if err != nil {
		return nil, WrapError(err, ErrJSONParsing).WithDetail("function", name)
	}
	var paramsMap map[string]any
	if err := json.Unmarshal(paramsJSON, &paramsMap); err != nil {
		return nil, WrapError(err, ErrJSONParsing).WithDetail("function", name)
	}
	return paramsMap, nil
}

func convertToResponsesToolChoice(toolChoice any) responses.ResponseNewParamsToolChoiceUnion {
	if strChoice, ok := toolChoice.(string); ok {
		switch strChoice {
		case "auto":
			return responses.ResponseNewParamsToolChoiceUnion{
				OfToolChoiceMode: param.NewOpt(responses.ToolChoiceOptionsAuto),
			}
		case "none":
			return responses.ResponseNewParamsToolChoiceUnion{
				OfToolChoiceMode: param.NewOpt(responses.ToolChoiceOptionsNone),
			}
		case "required":
			return responses.ResponseNewParamsToolChoiceUnion{
				OfToolChoiceMode: param.NewOpt(responses.ToolChoiceOptionsRequired),
			}
		}
	}
	return responses.ResponseNewParamsToolChoiceUnion{
		OfToolChoiceMode: param.NewOpt(responses.ToolChoiceOptionsAuto),
	}
}

func convertToReasoningEffort(effort string) shared.ReasoningEffort {
	switch strings.ToLower(effort) {
	case "low":
		return shared.ReasoningEffortLow
	case "medium":
		return shared.ReasoningEffortMedium
	case "high":
		return shared.ReasoningEffortHigh
	case "minimal":
		return shared.ReasoningEffortMinimal
	case "xhigh":
		return shared.ReasoningEffortXhigh
	default:
		return shared.ReasoningEffortMedium
	}
}

func convertToResponsesFormatParam(format *llm.ResponseFormat) (responses.ResponseFormatTextConfigUnionParam, error) {
	switch format.Type {
	case llm.JSONObject:
		return responses.ResponseFormatTextConfigUnionParam{
			OfJSONObject: &shared.ResponseFormatJSONObjectParam{},
		}, nil
	case llm.JSONSchema:
		schema, ok := format.JSONSchema.(map[string]any)
		if !ok {
			schemaBytes, err := json.Marshal(format.JSONSchema)
			if err != nil {
				return responses.ResponseFormatTextConfigUnionParam{}, WrapError(err, ErrJSONParsing)
			}
			var schemaMap map[string]any
			if err := json.Unmarshal(schemaBytes, &schemaMap); err != nil {
				return responses.ResponseFormatTextConfigUnionParam{}, WrapError(err, ErrJSONParsing)
			}
			schema = schemaMap
		}
		return responses.ResponseFormatTextConfigUnionParam{
			OfJSONSchema: &responses.ResponseFormatTextJSONSchemaConfigParam{
				Name:   "schema",
				Schema: schema,
			},
		}, nil
	default:
		return responses.ResponseFormatTextConfigUnionParam{
			OfText: &shared.ResponseFormatTextParam{},
		}, nil
	}
}

func convertFromResponsesResponse(resp *responses.Response) (llm.Response, error) {
	message := llm.Message{
		Role:    llm.RoleAssistant,
		Content: resp.OutputText(),
	}
	for _, item := range resp.Output {
		if item.Type == "function_call" {
			fc := item.AsFunctionCall()
			message.ToolCalls = append(message.ToolCalls, llm.ToolCall{
				ID:   fc.CallID,
				Type: "function",
				Function: llm.FunctionCall{Name: fc.Name, Arguments: fc.Arguments},
			})
		}
	}
	usage := llm.Usage{
		PromptTokens:     int(resp.Usage.InputTokens),
		CompletionTokens: int(resp.Usage.OutputTokens),
		TotalTokens:      int(resp.Usage.InputTokens + resp.Usage.OutputTokens),
	}
	return llm.Response{Message: message, Usage: usage}, nil
}

func convertToResponsesContentParts(parts []llm.ContentPart) (responses.ResponseInputMessageContentListParam, error) {
	result := make(responses.ResponseInputMessageContentListParam, 0, len(parts))
	for _, part := range parts {
		switch part.Type {
		case llm.ContentPartTypeText:
			result = append(result, responses.ResponseInputContentUnionParam{
				OfInputText: &responses.ResponseInputTextParam{Text: part.Text},
			})
		case llm.ContentPartTypeImageURL:
			if part.ImageURL == nil {
				return nil, errorRegistry.New(ErrInvalidMessage).WithDetail("error", "image_url content part missing image_url")
			}
			detail := responses.ResponseInputImageDetailAuto
			switch part.ImageURL.Detail {
			case llm.ImageDetailLow:
				detail = responses.ResponseInputImageDetailLow
			case llm.ImageDetailHigh:
				detail = responses.ResponseInputImageDetailHigh
			}
			result = append(result, responses.ResponseInputContentUnionParam{
				OfInputImage: &responses.ResponseInputImageParam{
					Detail:   detail,
					ImageURL: openai.String(part.ImageURL.URL),
				},
			})
		case llm.ContentPartTypeFile:
			if part.File == nil {
				return nil, errorRegistry.New(ErrInvalidMessage).WithDetail("error", "file content part missing file")
			}
			fileParam := &responses.ResponseInputFileParam{}
			if part.File.FileID != "" {
				fileParam.FileID = openai.String(part.File.FileID)
			}
			if part.File.FileData != "" {
				fileParam.FileData = openai.String(part.File.FileData)
			}
			if part.File.Filename != "" {
				fileParam.Filename = openai.String(part.File.Filename)
			}
			result = append(result, responses.ResponseInputContentUnionParam{OfInputFile: fileParam})
		case llm.ContentPartTypeInputAudio:
			continue // Responses API does not support audio input content parts
		default:
			return nil, errorRegistry.New(ErrInvalidMessage).
				WithDetail("error", fmt.Sprintf("unsupported content part type: %s", part.Type))
		}
	}
	return result, nil
}

func convertToFloat32Slice(input []float64) []float32 {
	result := make([]float32, len(input))
	for i, v := range input {
		result[i] = float32(v)
	}
	return result
}
