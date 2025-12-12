package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MrJeffLarry/redmine-cli/internal/config"
	expect "github.com/Netflix/go-expect"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

type testConfig struct {
	R *config.Red_t
	U *user
	M *http.ServeMux
	S *httptest.Server
	t *testing.T
	c *expect.Console
}

func common(t *testing.T) *testConfig {
	tc := &testConfig{}
	tc.R = &config.Red_t{
		Spinner: spinner.New(spinner.CharSets[9], 100*time.Millisecond),
	}
	tc.U = &user{
		User: userInfo{
			ID:          1,
			Login:       "test",
			Admin:       false,
			FirstName:   "test",
			LastName:    "test",
			AvatarUrl:   "test",
			TwofaScheme: "test",
			ApiKey:      "0987654321",
		},
	}

	tc.M = http.NewServeMux()
	tc.S = httptest.NewServer(tc.M)
	tc.R.Client = tc.S.Client()
	tc.t = t

	return tc
}

func (tc *testConfig) ExpectString(s string) {
	if _, err := tc.c.ExpectString(s); err != nil {
		tc.t.Helper()
		tc.t.Fatalf("ExpectString(%q) = %v", s, err)
	}
}

func (tc *testConfig) SendLine(s string) {
	if _, err := tc.c.SendLine(s); err != nil {
		tc.t.Helper()
		tc.t.Fatalf("SendLine(%q) = %v", s, err)
	}
}

func TestAuthLoginApiKeyNoServer(t *testing.T) {
	r := config.InitConfig()
	defer r.Save()
	tc := common(t)
	defer tc.S.Close()

	tc.M.HandleFunc("/users/current.json", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.Write([]byte(`"bad":parse"}`))
	})

	if loginApiKey(tc.R, &cobra.Command{}, tc.S.URL, tc.U.User.ApiKey) {
		t.Error("Wanted no server but got success")
	}

	if tc.U.User.ApiKey == tc.R.Config.ApiKey {
		t.Error("Wanted no ApiKey but got match")
	}
}

func TestAuthLoginApiKeyBadServer(t *testing.T) {
	r := config.InitConfig()
	defer r.Save()
	tc := common(t)
	defer tc.S.Close()

	if loginApiKey(tc.R, &cobra.Command{}, "", tc.U.User.ApiKey) {
		t.Error("Wanted bad server but got success")
	}

	if tc.U.User.ApiKey == tc.R.Config.ApiKey {
		t.Error("Wanted bad ApiKey but got match")
	}
}

func TestAuthLoginApiKeyBadCred(t *testing.T) {
	r := config.InitConfig()
	defer r.Save()
	tc := common(t)
	defer tc.S.Close()

	tc.M.HandleFunc("/users/current.json", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Invalid credentials"}`))
	})

	if loginApiKey(tc.R, &cobra.Command{}, tc.S.URL, tc.U.User.ApiKey) {
		t.Error("Wanted bad login but got success")
	}

	if tc.U.User.ApiKey == tc.R.Config.ApiKey {
		t.Error("Wanted bad ApiKey but got match")
	}
}

func TestAuthLoginApiKeyOk(t *testing.T) {
	r := config.InitConfig()
	defer r.Save()
	tc := common(t)
	defer tc.S.Close()

	tc.M.HandleFunc("/users/current.json", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			t.Errorf("Wanted Method[GET] got[%s]", r.Method)
			return
		}

		if r.Header.Get("X-Redmine-API-Key") != tc.U.User.ApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			t.Errorf("Wanted X-Redmine-API-Key[%s] got[%s]", tc.U.User.ApiKey, r.Header.Get("X-Redmine-API-Key"))
			return
		}

		body, err := json.Marshal(tc.U)
		if err != nil {
			t.Errorf("Wanted err[<nil>] got[%s]", err.Error())
		}
		w.Write(body)
	})

	loginApiKey(tc.R, &cobra.Command{}, tc.S.URL, tc.U.User.ApiKey)

	if tc.U.User.ApiKey != tc.R.Config.ApiKey {
		t.Errorf("Wanted ApiKey[%s] got[%s]", tc.U.User.ApiKey, tc.R.Config.ApiKey)
	}

	if tc.S.URL != tc.R.Config.Server {
		t.Errorf("Wanted Server[%s] got[%s]", tc.S.URL, tc.R.Config.Server)
	}
}

func TestAuthLoginApiKeyWithRID(t *testing.T) {
	r := config.InitConfig()
	defer r.Save()
	tc := common(t)
	defer tc.S.Close()

	tc.M.HandleFunc("/users/current.json", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			t.Errorf("Wanted Method[GET] got[%s]", r.Method)
			return
		}

		if r.Header.Get("X-Redmine-API-Key") != tc.U.User.ApiKey {
			w.WriteHeader(http.StatusUnauthorized)
			t.Errorf("Wanted X-Redmine-API-Key[%s] got[%s]", tc.U.User.ApiKey, r.Header.Get("X-Redmine-API-Key"))
			return
		}

		body, err := json.Marshal(tc.U)
		if err != nil {
			t.Errorf("Wanted err[<nil>] got[%s]", err.Error())
		}
		w.Write(body)
	})

	// Create a command with RID flag
	cmd := &cobra.Command{}
	cmd.Flags().String(config.RID_FLAG, "", "Redmine instance ID")
	cmd.Flags().Set(config.RID_FLAG, "2")

	loginApiKey(tc.R, cmd, tc.S.URL, tc.U.User.ApiKey)

	// Verify multi-mode is enabled
	if !tc.R.UseMultiMode {
		t.Error("Expected UseMultiMode to be true when using --rid flag")
	}

	// Verify RID is set
	if tc.R.RID != "2" {
		t.Errorf("Wanted RID[2] got[%s]", tc.R.RID)
	}

	// Verify config is correct
	if tc.U.User.ApiKey != tc.R.Config.ApiKey {
		t.Errorf("Wanted ApiKey[%s] got[%s]", tc.U.User.ApiKey, tc.R.Config.ApiKey)
	}

	if tc.S.URL != tc.R.Config.Server {
		t.Errorf("Wanted Server[%s] got[%s]", tc.S.URL, tc.R.Config.Server)
	}
}
