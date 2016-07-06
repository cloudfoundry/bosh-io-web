package controllers

import (
	"errors"
	"net/http"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhchecks "github.com/cppforlife/bosh-hub/checksumsrepo"
)

type ChecksumReqMatch struct {
	Token    string
	Includes []string
}

func (m ChecksumReqMatch) Validate() error {
	if len(m.Token) < 40 {
		return errors.New("Expected token to be >40 chars long")
	}

	if len(m.Includes) == 0 {
		return errors.New("Expected to have at least one include directive")
	}

	for _, inc := range m.Includes {
		if len(inc) == 0 {
			return errors.New("Expected to include directive to be non-empty")
		}
	}

	return nil
}

type ChecksumsController struct {
	allowedMatches []ChecksumReqMatch
	checksumsRepo  bhchecks.ChecksumsRepository
	logger         boshlog.Logger
}

func NewChecksumsController(
	allowedMatches []ChecksumReqMatch,
	checksumsRepo bhchecks.ChecksumsRepository,
	logger boshlog.Logger,
) ChecksumsController {
	return ChecksumsController{
		allowedMatches: allowedMatches,
		checksumsRepo:  checksumsRepo,
		logger:         logger,
	}
}

func (c ChecksumsController) Find(req *http.Request, r martrend.Render, params mart.Params) {
	key, err := c.authorize(req, params)
	if err != nil {
		r.JSON(401, map[string]string{"error": "Unauthorized: " + err.Error()})
		return
	}

	rec, err := c.checksumsRepo.Find(key)
	if err != nil {
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, map[string]string{"sha1": rec.SHA1})
}

func (c ChecksumsController) Save(req *http.Request, r martrend.Render, params mart.Params) {
	key, err := c.authorize(req, params)
	if err != nil {
		r.JSON(401, map[string]string{"error": "Unauthorized: " + err.Error()})
		return
	}

	rec := bhchecks.ChecksumRec{SHA1: req.FormValue("sha1")}

	err = c.checksumsRepo.Save(key, rec)
	if err != nil {
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	r.JSON(200, map[string]string{"sha1": rec.SHA1})
}

func (c ChecksumsController) authorize(req *http.Request, params mart.Params) (string, error) {
	accessedKey := strings.TrimSpace(params["_1"])

	if len(accessedKey) == 0 {
		return "", errors.New("Key must be non-empty")
	}

	// Presented in the header as "Bearer blabhlbah"
	authHeader := strings.ToLower(req.Header.Get("Authorization"))
	givenToken := strings.TrimSpace(strings.TrimPrefix(authHeader, "bearer "))

	if len(givenToken) == 0 {
		return "", errors.New("Token must be non-empty")
	}

	var foundMatches []ChecksumReqMatch

	for _, priv := range c.allowedMatches {
		for _, inc := range priv.Includes {
			if strings.Contains(accessedKey, inc) {
				foundMatches = append(foundMatches, priv)
			}
		}
	}

	if len(foundMatches) != 1 {
		return "", errors.New("Expected to match exactly one rule")
	}

	if foundMatches[0].Token != givenToken {
		return "", errors.New("Token mismatch")
	}

	return accessedKey, nil
}
