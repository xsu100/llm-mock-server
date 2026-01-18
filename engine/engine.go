package engine

import (
	"math/rand"
	"regexp"
	"strings"

	"llm-mock-server/models"
)

type Engine struct {
	rules []models.MockRule
}

func NewEngine(rules []models.MockRule) *Engine {
	return &Engine{rules: rules}
}

func (e *Engine) MatchChat(req models.ChatCompletionRequest) (string, bool) {
	// Combine all messages to search for patterns
	var fullText strings.Builder
	for _, m := range req.Messages {
		fullText.WriteString(m.Content)
		fullText.WriteString(" ")
	}
	content := fullText.String()

	// 1. Try Regex/Keyword Match
	for _, rule := range e.rules {
		if rule.Type == models.RuleTypeRegex {
			matched, _ := regexp.MatchString(rule.Pattern, content)
			if matched {
				return rule.Response, true
			}
		}
	}

	// 2. Try Fixed Match
	for _, rule := range e.rules {
		if rule.Type == models.RuleTypeFixed {
			if strings.Contains(content, rule.Pattern) {
				return rule.Response, true
			}
		}
	}

	// 3. Try Random (if any random rules exist)
	var randomRules []models.MockRule
	for _, rule := range e.rules {
		if rule.Type == models.RuleTypeRandom {
			randomRules = append(randomRules, rule)
		}
	}

	if len(randomRules) > 0 {
		idx := rand.Intn(len(randomRules))
		return randomRules[idx].Response, true
	}

	return "", false
}

func (e *Engine) MatchImage(req models.ImageGenerationRequest) (string, bool) {
	content := req.Prompt

	for _, rule := range e.rules {
		if rule.Type == models.RuleTypeRegex {
			matched, _ := regexp.MatchString(rule.Pattern, content)
			if matched {
				return rule.Response, true
			}
		}
		if rule.Type == models.RuleTypeFixed {
			if strings.Contains(content, rule.Pattern) {
				return rule.Response, true
			}
		}
	}

	return "", false
}
