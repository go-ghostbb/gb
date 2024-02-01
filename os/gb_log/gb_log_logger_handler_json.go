package gblog

import (
	"context"
	"ghostbb.io/internal/json"
	gbconv "ghostbb.io/util/gb_conv"
)

// HandlerOutputJson is the structure outputting logging content as single json.
type HandlerOutputJson struct {
	Time       string `json:""`           // Formatted time string, like "2016-01-09 12:00:00".
	TraceId    string `json:",omitempty"` // Trace id, only available if tracing is enabled.
	CtxStr     string `json:",omitempty"` // The retrieved context value string from context, only available if Config.CtxKeys configured.
	Level      string `json:""`           // Formatted level string, like "DEBU", "ERRO", etc. Eg: ERRO
	CallerPath string `json:",omitempty"` // The source file path and its line number that calls logging, only available if F_FILE_SHORT or F_FILE_LONG set.
	CallerFunc string `json:",omitempty"` // The source function name that calls logging, only available if F_CALLER_FN set.
	Prefix     string `json:",omitempty"` // Custom prefix string for logging content.
	Content    string `json:""`           // Content is the main logging content, containing error stack string produced by logger.
	Stack      string `json:",omitempty"` // Stack string produced by logger, only available if Config.StStatus configured.
}

// HandlerJson is a handler for output logging content as a single json string.
func HandlerJson(ctx context.Context, in *HandlerInput) {
	output := HandlerOutputJson{
		Time:       in.TimeFormat,
		TraceId:    in.TraceId,
		CtxStr:     in.CtxStr,
		Level:      in.LevelFormat,
		CallerFunc: in.CallerFunc,
		CallerPath: in.CallerPath,
		Prefix:     in.Prefix,
		Content:    in.Content,
		Stack:      in.Stack,
	}
	// Convert values string content.
	var valueContent string
	for _, v := range in.Values {
		valueContent = gbconv.String(v)
		if len(valueContent) == 0 {
			continue
		}
		if len(output.Content) > 0 {
			if output.Content[len(output.Content)-1] == '\n' {
				// Remove one blank line(\n\n).
				if valueContent[0] == '\n' {
					valueContent = valueContent[1:]
				}
				output.Content += valueContent
			} else {
				output.Content += " " + valueContent
			}
		} else {
			output.Content += valueContent
		}
	}
	// Output json content.
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	in.Buffer.Write(jsonBytes)
	in.Buffer.Write([]byte("\n"))
	in.Next(ctx)
}
