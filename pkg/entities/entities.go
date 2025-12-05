package entities

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/port-experimental/port-go-sdk/pkg/porter"
)

// Service handles entity endpoints.
type Service struct {
	doer Doer
}

// New creates an entity service.
func New(doer Doer) *Service {
	return &Service{doer: doer}
}

// Doer matches client.Client for dependency injection.
type Doer interface {
	Do(ctx context.Context, method, path string, body any, out any) error
}

// Entity represents a Port entity.
type Entity struct {
	Identifier string                 `json:"identifier"`
	Blueprint  string                 `json:"blueprint"`
	Title      string                 `json:"title,omitempty"`
	Icon       string                 `json:"icon,omitempty"`
	Team       string                 `json:"team,omitempty"`
	Properties map[string]any         `json:"properties,omitempty"`
	Relations  map[string][]string    `json:"relations,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// ListOptions control pagination/filtering.
type ListOptions struct {
	Include []string
	Exclude []string

	ExcludeCalculatedProperties             bool
	AttachTitleToRelation                   bool
	AttachIdentifierToTitleMirrorProperties bool
	AllowPartialResults                     bool

	// Deprecated: Search-only controls retained for compatibility. Use Search/SearchBlueprint instead.
	Query map[string]any
	From  string
	Limit int
}

// ListResponse wraps entity lists returned from the search endpoint.
type ListResponse struct {
	Entities []Entity `json:"entities"`
	Next     string   `json:"next,omitempty"`
	OK       bool     `json:"ok"`
}

// Create creates a new entity.
func (s *Service) Create(ctx context.Context, blueprint string, ent Entity) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities", url.PathEscape(blueprint))
	return s.doer.Do(ctx, "POST", path, entityPayload(ent), nil)
}

// Upsert creates or updates an entity.
func (s *Service) Upsert(ctx context.Context, blueprint string, ent Entity) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities?upsert=true&merge=true", url.PathEscape(blueprint))
	return s.doer.Do(ctx, "POST", path, entityPayload(ent), nil)
}

// Get fetches an entity by identifier.
func (s *Service) Get(ctx context.Context, blueprint, identifier string) (Entity, error) {
	var out Entity
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s", url.PathEscape(blueprint), url.PathEscape(identifier))
	err := s.doer.Do(ctx, "GET", path, nil, &out)
	return out, err
}

// Delete removes an entity.
func (s *Service) Delete(ctx context.Context, blueprint, identifier string) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s", url.PathEscape(blueprint), url.PathEscape(identifier))
	return s.doer.Do(ctx, "DELETE", path, nil, nil)
}

// List returns entities for a blueprint with optional filters.
func (s *Service) List(ctx context.Context, blueprint string, opts *ListOptions) (ListResponse, error) {
	var out ListResponse
	path := fmt.Sprintf("/v1/blueprints/%s/entities", url.PathEscape(blueprint))
	if opts != nil {
		if opts.Query != nil || opts.From != "" || opts.Limit > 0 {
			return ListResponse{}, fmt.Errorf("entities list: Query/From/Limit are not supported, use Search/SearchBlueprint instead")
		}
		values := url.Values{}
		for _, inc := range opts.Include {
			values.Add("include", inc)
		}
		for _, exc := range opts.Exclude {
			values.Add("exclude", exc)
		}
		if opts.ExcludeCalculatedProperties {
			values.Set("exclude_calculated_properties", strconv.FormatBool(true))
		}
		if opts.AttachTitleToRelation {
			values.Set("attach_title_to_relation", strconv.FormatBool(true))
		}
		if opts.AttachIdentifierToTitleMirrorProperties {
			values.Set("attach_identifier_to_title_mirror_properties", strconv.FormatBool(true))
		}
		if opts.AllowPartialResults {
			values.Set("allow_partial_results", strconv.FormatBool(true))
		}
		if qs := values.Encode(); qs != "" {
			path += "?" + qs
		}
	}
	err := s.doer.Do(ctx, "GET", path, nil, &out)
	return out, err
}

// Update applies a partial update to entity properties (merge=true).
func (s *Service) Update(ctx context.Context, blueprint, identifier string, properties map[string]any) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s", url.PathEscape(blueprint), url.PathEscape(identifier))
	payload := map[string]any{
		"properties": properties,
	}
	if err := s.doer.Do(ctx, "PATCH", path, payload, nil); err != nil {
		var perr *porter.Error
		if errors.As(err, &perr) && perr.StatusCode == 422 {
			return s.doer.Do(ctx, "PUT", path, payload, nil)
		}
		return err
	}
	return nil
}

// LinkRelation links targets to a relation.
func (s *Service) LinkRelation(ctx context.Context, blueprint, identifier, relation string, targets []string) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s/relations/%s", url.PathEscape(blueprint), url.PathEscape(identifier), url.PathEscape(relation))
	payload := map[string]any{"identifiers": targets}
	return s.doer.Do(ctx, "POST", path, payload, nil)
}

// UnlinkRelation removes relation links.
func (s *Service) UnlinkRelation(ctx context.Context, blueprint, identifier, relation string, targets []string) error {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/%s/relations/%s", url.PathEscape(blueprint), url.PathEscape(identifier), url.PathEscape(relation))
	payload := map[string]any{"identifiers": targets}
	return s.doer.Do(ctx, "DELETE", path, payload, nil)
}

// BulkEntitiesResponse describes the outcome of bulk create/upsert.
type BulkEntitiesResponse struct {
	OK       bool               `json:"ok"`
	Entities []BulkEntityStatus `json:"entities"`
}

// BulkEntityStatus reports the per-entity result in a bulk create call.
type BulkEntityStatus struct {
	Created        bool                   `json:"created"`
	Identifier     string                 `json:"identifier"`
	Index          int                    `json:"index"`
	AdditionalData map[string]interface{} `json:"additionalData,omitempty"`
}

// BulkUpsert creates or updates up to 20 entities in a single call.
func (s *Service) BulkUpsert(ctx context.Context, blueprint string, entities []Entity) (BulkEntitiesResponse, error) {
	if len(entities) == 0 {
		return BulkEntitiesResponse{}, fmt.Errorf("entities: at least one entity required for bulk upsert")
	}
	items := make([]map[string]any, len(entities))
	for i, ent := range entities {
		items[i] = entityPayload(ent)
	}
	payload := map[string]any{"entities": items}
	path := fmt.Sprintf("/v1/blueprints/%s/entities/bulk", url.PathEscape(blueprint))
	var resp BulkEntitiesResponse
	err := s.doer.Do(ctx, "POST", path, payload, &resp)
	return resp, err
}

// BulkDeleteOptions customize bulk deletion.
type BulkDeleteOptions struct {
	DeleteDependents bool
	RunID            string
}

// BulkDeleteResponse lists the identifiers removed via bulk delete.
type BulkDeleteResponse struct {
	OK              bool     `json:"ok"`
	DeletedEntities []string `json:"deletedEntities"`
}

// BulkDelete removes up to 100 entities from a blueprint.
func (s *Service) BulkDelete(ctx context.Context, blueprint string, identifiers []string, opts *BulkDeleteOptions) (BulkDeleteResponse, error) {
	if len(identifiers) == 0 {
		return BulkDeleteResponse{}, fmt.Errorf("entities: at least one identifier required for bulk delete")
	}
	options := BulkDeleteOptions{}
	if opts != nil {
		options = *opts
	}
	values := url.Values{}
	values.Set("delete_dependents", strconv.FormatBool(options.DeleteDependents))
	if options.RunID != "" {
		values.Set("run_id", options.RunID)
	}
	path := fmt.Sprintf("/v1/blueprints/%s/bulk/entities/delete", url.PathEscape(blueprint))
	if qs := values.Encode(); qs != "" {
		path += "?" + qs
	}
	body := map[string]any{"entities": identifiers}
	var resp BulkDeleteResponse
	err := s.doer.Do(ctx, "POST", path, body, &resp)
	return resp, err
}

// SearchOptions control the /entities/search POST body.
type SearchOptions struct {
	Query   map[string]any
	Include []string
	Exclude []string
	From    string
	Limit   int
}

// Search runs a cross-blueprint entities search.
func (s *Service) Search(ctx context.Context, opts SearchOptions) (ListResponse, error) {
	return s.search(ctx, "/v1/entities/search", opts)
}

// SearchBlueprint searches within a given blueprint.
func (s *Service) SearchBlueprint(ctx context.Context, blueprint string, opts SearchOptions) (ListResponse, error) {
	path := fmt.Sprintf("/v1/blueprints/%s/entities/search", url.PathEscape(blueprint))
	return s.search(ctx, path, opts)
}

func (s *Service) search(ctx context.Context, path string, opts SearchOptions) (ListResponse, error) {
	body := map[string]any{}
	if len(opts.Include) > 0 {
		body["include"] = opts.Include
	}
	if len(opts.Exclude) > 0 {
		body["exclude"] = opts.Exclude
	}
	if opts.Query != nil {
		body["query"] = opts.Query
	}
	if opts.From != "" {
		body["from"] = opts.From
	}
	if opts.Limit > 0 {
		body["limit"] = opts.Limit
	}
	if len(body) == 0 {
		body = nil
	}
	var out ListResponse
	err := s.doer.Do(ctx, "POST", path, body, &out)
	return out, err
}

func entityPayload(ent Entity) map[string]any {
	payload := map[string]any{
		"identifier": ent.Identifier,
		"properties": ent.Properties,
	}
	if ent.Title != "" {
		payload["title"] = ent.Title
	}
	if ent.Icon != "" {
		payload["icon"] = ent.Icon
	}
	if ent.Team != "" {
		payload["team"] = ent.Team
	}
	if rel := cloneRelations(ent.Relations); rel != nil {
		payload["relations"] = rel
	}
	return payload
}

func cloneRelations(in map[string][]string) map[string]any {
	if len(in) == 0 {
		return nil
	}
	rel := make(map[string]any, len(in))
	for k, v := range in {
		dst := make([]string, len(v))
		copy(dst, v)
		rel[k] = dst
	}
	return rel
}

// AggregateRequest mirrors the flexible aggregate payload.
type AggregateRequest map[string]any

// AggregateResponse represents /v1/entities/aggregate output.
type AggregateResponse struct {
	OK                 bool     `json:"ok"`
	Entities           []Entity `json:"entities"`
	MatchingBlueprints []string `json:"matchingBlueprints,omitempty"`
	FailedBlueprints   []string `json:"failedBlueprints,omitempty"`
}

// Aggregate executes /v1/entities/aggregate with the provided payload.
func (s *Service) Aggregate(ctx context.Context, req AggregateRequest) (AggregateResponse, error) {
	if len(req) == 0 {
		return AggregateResponse{}, fmt.Errorf("entities: aggregate request cannot be empty")
	}
	var resp AggregateResponse
	err := s.doer.Do(ctx, "POST", "/v1/entities/aggregate", req, &resp)
	return resp, err
}

// AggregateOverTimeRequest configures /v1/entities/aggregate-over-time.
type AggregateOverTimeRequest struct {
	Blueprint         string             `json:"blueprint"`
	TimeRange         AggregateTimeRange `json:"timeRange"`
	TimeInterval      string             `json:"timeInterval"`
	Query             map[string]any     `json:"query"`
	MeasureTimeBy     string             `json:"measureTimeBy"`
	AggregationType   string             `json:"aggregationType"`
	Func              string             `json:"func"`
	Properties        []string           `json:"properties,omitempty"`
	BreakdownProperty string             `json:"breakdownProperty,omitempty"`
}

// AggregateTimeRange defines preset + timezone.
type AggregateTimeRange struct {
	Preset   string `json:"preset"`
	TimeZone string `json:"timeZone,omitempty"`
}

// AggregateOverTimeResponse captures the returned series.
type AggregateOverTimeResponse struct {
	OK     bool                    `json:"ok"`
	Result AggregateOverTimeResult `json:"result"`
}

// AggregateOverTimeResult includes the time bounds and rows.
type AggregateOverTimeResult struct {
	MinDate float64              `json:"minDate"`
	MaxDate float64              `json:"maxDate"`
	Data    []map[string]float64 `json:"data"`
}

// AggregateOverTime calls /v1/entities/aggregate-over-time.
func (s *Service) AggregateOverTime(ctx context.Context, req AggregateOverTimeRequest) (AggregateOverTimeResponse, error) {
	var resp AggregateOverTimeResponse
	err := s.doer.Do(ctx, "POST", "/v1/entities/aggregate-over-time", req, &resp)
	return resp, err
}

// PropertiesHistoryRequest configures /v1/entities/properties-history.
type PropertiesHistoryRequest struct {
	EntityIdentifier    string                      `json:"entityIdentifier"`
	BlueprintIdentifier string                      `json:"blueprintIdentifier"`
	PropertyNames       []string                    `json:"propertyNames"`
	TimeInterval        string                      `json:"timeInterval,omitempty"`
	TimeRange           *PropertiesHistoryTimeRange `json:"timeRange,omitempty"`
}

// PropertiesHistoryTimeRange describes the preset and zone.
type PropertiesHistoryTimeRange struct {
	Preset   string `json:"preset"`
	TimeZone string `json:"timeZone,omitempty"`
}

// PropertiesHistoryResponse contains time-series details.
type PropertiesHistoryResponse struct {
	OK     bool                    `json:"ok"`
	Result PropertiesHistoryResult `json:"result"`
}

// PropertiesHistoryResult holds the numeric samples.
type PropertiesHistoryResult struct {
	MinDate float64   `json:"minDate"`
	MaxDate float64   `json:"maxDate"`
	Data    []float64 `json:"data"`
}

// PropertiesHistory hits /v1/entities/properties-history.
func (s *Service) PropertiesHistory(ctx context.Context, req PropertiesHistoryRequest) (PropertiesHistoryResponse, error) {
	var resp PropertiesHistoryResponse
	err := s.doer.Do(ctx, "POST", "/v1/entities/properties-history", req, &resp)
	return resp, err
}
