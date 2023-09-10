package main

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

func start(submissionID uint) {
	var submission Submission
	d.First(&submission, submissionID)
	defer d.Save(&submission)

	path, err := os.MkdirTemp("", "sandbox-*")
	if err != nil {
		submission.Verdict = "Internal Error"
		return
	}
	defer os.RemoveAll(path)

	err = os.WriteFile(path+"/"+submission.FileName, []byte(submission.Content), os.ModePerm)
	if err != nil {
		submission.Verdict = "Internal Error"
		return
	}

	var problem Problem
	d.First(&problem, submission.ProblemID)

	if len(problem.InputName) > 0 {
		err = os.WriteFile(path+"/"+problem.InputName, []byte(problem.Input), os.ModePerm)
		if err != nil {
			submission.Verdict = "Internal Error"
			return
		}
	}

	compileStatus := compile(path + "/" + submission.FileName)
	if !compileStatus {
		submission.Verdict = "Compile Time Error"
		return
	}

	output, runtimeStatus := run(submission.FileName, path)
	if runtimeStatus != "Runtime Finished" {
		submission.Verdict = runtimeStatus
		return
	}

	submission.Verdict = judge(format(strings.Split(string(output), "\n")), format(strings.Split(problem.Output, "\n")))
}

func judge(codeOutputList []string, sampleOutputList []string) string {
	if len(codeOutputList) != len(sampleOutputList) {
		return "Wrong Answer"
	}

	for i := 0; i < len(codeOutputList); i++ {
		if codeOutputList[i] != sampleOutputList[i] {
			return "Wrong Answer"
		}
	}
	return "Correct Answer"
}

func format(outputList []string) []string {
	for {
		if len(outputList) == 0 || !whitespace(outputList[0]) {
			break
		}

		outputList = outputList[1:]
	}

	for {
		if len(outputList) == 0 || !whitespace(outputList[len(outputList)-1]) {
			break
		}
		outputList = outputList[:len(outputList)-1]
	}

	for index, line := range outputList {
		outputList[index] = strings.TrimRightFunc(line, func(r rune) bool {
			return unicode.IsSpace(r)
		})
	}

	return outputList
}

func whitespace(str string) bool {
	for _, c := range str {
		if !unicode.IsSpace(c) {
			return false
		}
	}
	return true
}

func compile(path string) bool {
	_, err := exec.Command(
		"javac",
		path,
		"-Xlint:unchecked").CombinedOutput()
	return err == nil
}

func run(fileName string, path string) ([]byte, string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx,
		"java",
		"-Djava.security.manager",
		string(fileName[0:strings.LastIndex(fileName, ".")]))

	cmd.Dir = path
	out, err := cmd.CombinedOutput()

	index := 0
	count := 2
	for {
		if index == len(out) || count == 0 {
			break
		}

		if out[index] == '\n' {
			count--
		}
		index++
	}
	out = out[index:]

	if ctx.Err() == context.DeadlineExceeded {
		return out, "Time Limit Exceeded"
	} else {
		if err == nil {
			return out, "Runtime Finished"
		} else {
			return out, "Runtime Error"
		}
	}
}
