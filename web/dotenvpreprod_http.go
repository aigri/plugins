package web

import (
        "github.com/LeakIX/l9format"
        "regexp"
        "strings"
)

type EnvPreprodHttpPlugin struct {
        l9format.ServicePluginBase
}

func (EnvPreprodHttpPlugin) GetVersion() (int, int, int) {
        return 0, 0, 2
}

func (EnvPreprodHttpPlugin) GetRequests() []l9format.WebPluginRequest {
        return []l9format.WebPluginRequest{{
                Method:  "GET",
                Path:    "/.env.preprod",
                Headers: map[string]string{},
                Body:    []byte(""),
        }}
}

func (EnvPreprodHttpPlugin) GetName() string {
        return "EnvPreprodHttpPlugin"
}

func (EnvPreprodHttpPlugin) GetStage() string {
        return "open"
}

func (plugin EnvPreprodHttpPlugin) Verify(request l9format.WebPluginRequest, response l9format.WebPluginResponse, event *l9format.L9Event, options map[string]string) (hasLeak bool) {
        if !request.EqualAny(plugin.GetRequests()) || response.Response.StatusCode != 200 {
                return false
        }
        lowerBody := strings.ToLower(string(response.Body))
        if len(lowerBody) < 10 {
                return false
        }

        // Regex pour détecter des informations sensibles
        regexPattern := `(?i)(app_env=|db_host=|\bAKIA[A-Z0-9]{16}\b|SG\.[a-zA-Z0-9_-]{22}\.[a-zA-Z0-9_-]{43}|smtp\.mailgun\.org|smtp-relay\.sendinblue\.com)`
        match, _ := regexp.MatchString(regexPattern, lowerBody)

        if match {
                event.Service.Software.Name = "EnvironmentFile"
                event.Leak.Type = "config_leak"
                event.Leak.Severity = "high"
                event.AddTag("potential-leak")
                event.Summary = "Found sensitive information in /.env.preprod:\n" + string(response.Body)
                return true
        }

        return false
}
