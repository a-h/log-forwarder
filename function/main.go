package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("failed to init logger: %w", err))
	}
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.CloudwatchLogsEvent) (err error) {
	data, err := e.AWSLogs.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse logs data: %v", err)
	}

	// Create a pipe to send the JSON to the stdin of LogStash.
	r, w := io.Pipe()

	// Start it running in the background.
	go func() {
		// Unmarshal the entries into JSON.
		for _, e := range data.LogEvents {
			// Unmarshal JSON.
			m := make(map[string]interface{})
			err = json.Unmarshal([]byte(e.Message), &m)
			if err != nil {
				log.Error("failed to unmarshal log event", zap.Error(err))
				continue
			}
			// Add fields.
			m["logGroup"] = data.LogGroup
			m["logStream"] = data.LogStream
			// Write to output.
			err = json.NewEncoder(w).Encode(m)
			if err != nil {
				log.Error("failed to remarshal modified log event", zap.Error(err))
				continue
			}
		}
		// Close off the pipe which will cause the stdin to be closed.
		w.Close()
	}()

	return run(r)
}

func run(src io.Reader) (err error) {
	cmd := exec.Command("/usr/share/logstash/bin/logstash")
	cmd.Dir, err = os.Getwd()
	if err != nil {
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	input, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get input pipe: %w", err)
	}
	go func() {
		_, err = io.Copy(input, src)
		if err != nil {
			return
		}
		input.Close()
	}()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}
	if exitCode := cmd.ProcessState.ExitCode(); exitCode != 0 {
		return fmt.Errorf("non-zero exit code: %v", exitCode)
	}
	return nil
}
