package collector

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/go-github/v57/github"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/oauth2"

	"github.com/l13t/github_rate_limit_exporter/internal/config"
)

// Collector collects GitHub API rate limit metrics
type Collector struct {
	users   []config.User
	clients map[string]*github.Client

	// Prometheus metrics
	coreLimit     *prometheus.GaugeVec
	coreRemaining *prometheus.GaugeVec
	coreUsed      *prometheus.GaugeVec
	coreReset     *prometheus.GaugeVec

	searchLimit     *prometheus.GaugeVec
	searchRemaining *prometheus.GaugeVec
	searchUsed      *prometheus.GaugeVec
	searchReset     *prometheus.GaugeVec

	graphqlLimit     *prometheus.GaugeVec
	graphqlRemaining *prometheus.GaugeVec
	graphqlUsed      *prometheus.GaugeVec
	graphqlReset     *prometheus.GaugeVec

	integrationManifestLimit     *prometheus.GaugeVec
	integrationManifestRemaining *prometheus.GaugeVec
	integrationManifestUsed      *prometheus.GaugeVec
	integrationManifestReset     *prometheus.GaugeVec

	mu sync.RWMutex
}

// NewCollector creates a new GitHub rate limit collector
func NewCollector(users []config.User) *Collector {
	c := &Collector{
		users:   users,
		clients: make(map[string]*github.Client),
	}

	// Initialize Prometheus metrics
	c.coreLimit = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_core_limit",
			Help: "GitHub API core rate limit",
		},
		[]string{"user"},
	)

	c.coreRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_core_remaining",
			Help: "GitHub API core rate limit remaining",
		},
		[]string{"user"},
	)

	c.coreUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_core_used",
			Help: "GitHub API core rate limit used",
		},
		[]string{"user"},
	)

	c.coreReset = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_core_reset_timestamp",
			Help: "GitHub API core rate limit reset timestamp",
		},
		[]string{"user"},
	)

	c.searchLimit = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_search_limit",
			Help: "GitHub API search rate limit",
		},
		[]string{"user"},
	)

	c.searchRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_search_remaining",
			Help: "GitHub API search rate limit remaining",
		},
		[]string{"user"},
	)

	c.searchUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_search_used",
			Help: "GitHub API search rate limit used",
		},
		[]string{"user"},
	)

	c.searchReset = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_search_reset_timestamp",
			Help: "GitHub API search rate limit reset timestamp",
		},
		[]string{"user"},
	)

	c.graphqlLimit = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_graphql_limit",
			Help: "GitHub API GraphQL rate limit",
		},
		[]string{"user"},
	)

	c.graphqlRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_graphql_remaining",
			Help: "GitHub API GraphQL rate limit remaining",
		},
		[]string{"user"},
	)

	c.graphqlUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_graphql_used",
			Help: "GitHub API GraphQL rate limit used",
		},
		[]string{"user"},
	)

	c.graphqlReset = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_graphql_reset_timestamp",
			Help: "GitHub API GraphQL rate limit reset timestamp",
		},
		[]string{"user"},
	)

	c.integrationManifestLimit = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_integration_manifest_limit",
			Help: "GitHub API integration manifest rate limit",
		},
		[]string{"user"},
	)

	c.integrationManifestRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_integration_manifest_remaining",
			Help: "GitHub API integration manifest rate limit remaining",
		},
		[]string{"user"},
	)

	c.integrationManifestUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_integration_manifest_used",
			Help: "GitHub API integration manifest rate limit used",
		},
		[]string{"user"},
	)

	c.integrationManifestReset = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "github_rate_limit_integration_manifest_reset_timestamp",
			Help: "GitHub API integration manifest rate limit reset timestamp",
		},
		[]string{"user"},
	)

	// Initialize GitHub clients for each user
	for _, user := range users {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: user.Token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		c.clients[user.Name] = github.NewClient(tc)
	}

	return c
}

// Describe implements prometheus.Collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.coreLimit.Describe(ch)
	c.coreRemaining.Describe(ch)
	c.coreUsed.Describe(ch)
	c.coreReset.Describe(ch)

	c.searchLimit.Describe(ch)
	c.searchRemaining.Describe(ch)
	c.searchUsed.Describe(ch)
	c.searchReset.Describe(ch)

	c.graphqlLimit.Describe(ch)
	c.graphqlRemaining.Describe(ch)
	c.graphqlUsed.Describe(ch)
	c.graphqlReset.Describe(ch)

	c.integrationManifestLimit.Describe(ch)
	c.integrationManifestRemaining.Describe(ch)
	c.integrationManifestUsed.Describe(ch)
	c.integrationManifestReset.Describe(ch)
}

// Collect implements prometheus.Collector
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.coreLimit.Collect(ch)
	c.coreRemaining.Collect(ch)
	c.coreUsed.Collect(ch)
	c.coreReset.Collect(ch)

	c.searchLimit.Collect(ch)
	c.searchRemaining.Collect(ch)
	c.searchUsed.Collect(ch)
	c.searchReset.Collect(ch)

	c.graphqlLimit.Collect(ch)
	c.graphqlRemaining.Collect(ch)
	c.graphqlUsed.Collect(ch)
	c.graphqlReset.Collect(ch)

	c.integrationManifestLimit.Collect(ch)
	c.integrationManifestRemaining.Collect(ch)
	c.integrationManifestUsed.Collect(ch)
	c.integrationManifestReset.Collect(ch)
}

// Update fetches the latest rate limit data from GitHub API
func (c *Collector) Update(ctx context.Context) {
	var wg sync.WaitGroup

	for _, user := range c.users {
		wg.Add(1)
		go func(u config.User) {
			defer wg.Done()
			c.updateUserRateLimits(ctx, u)
		}(user)
	}

	wg.Wait()
}

func (c *Collector) updateUserRateLimits(ctx context.Context, user config.User) {
	client, ok := c.clients[user.Name]
	if !ok {
		log.Printf("No client found for user %s", user.Name)
		return
	}

	rateLimits, _, err := client.RateLimit.Get(ctx)
	if err != nil {
		log.Printf("Error fetching rate limits for user %s: %v", user.Name, err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Core rate limits
	if rateLimits.Core != nil {
		c.coreLimit.WithLabelValues(user.Name).Set(float64(rateLimits.Core.Limit))
		c.coreRemaining.WithLabelValues(user.Name).Set(float64(rateLimits.Core.Remaining))
		c.coreUsed.WithLabelValues(user.Name).Set(float64(rateLimits.Core.Limit - rateLimits.Core.Remaining))
		c.coreReset.WithLabelValues(user.Name).Set(float64(rateLimits.Core.Reset.Unix()))
	}

	// Search rate limits
	if rateLimits.Search != nil {
		c.searchLimit.WithLabelValues(user.Name).Set(float64(rateLimits.Search.Limit))
		c.searchRemaining.WithLabelValues(user.Name).Set(float64(rateLimits.Search.Remaining))
		c.searchUsed.WithLabelValues(user.Name).Set(float64(rateLimits.Search.Limit - rateLimits.Search.Remaining))
		c.searchReset.WithLabelValues(user.Name).Set(float64(rateLimits.Search.Reset.Unix()))
	}

	// GraphQL rate limits
	if rateLimits.GraphQL != nil {
		c.graphqlLimit.WithLabelValues(user.Name).Set(float64(rateLimits.GraphQL.Limit))
		c.graphqlRemaining.WithLabelValues(user.Name).Set(float64(rateLimits.GraphQL.Remaining))
		c.graphqlUsed.WithLabelValues(user.Name).Set(float64(rateLimits.GraphQL.Limit - rateLimits.GraphQL.Remaining))
		c.graphqlReset.WithLabelValues(user.Name).Set(float64(rateLimits.GraphQL.Reset.Unix()))
	}

	// Integration Manifest rate limits
	if rateLimits.IntegrationManifest != nil {
		c.integrationManifestLimit.WithLabelValues(user.Name).Set(float64(rateLimits.IntegrationManifest.Limit))
		c.integrationManifestRemaining.WithLabelValues(user.Name).Set(float64(rateLimits.IntegrationManifest.Remaining))
		c.integrationManifestUsed.WithLabelValues(user.Name).Set(float64(rateLimits.IntegrationManifest.Limit - rateLimits.IntegrationManifest.Remaining))
		c.integrationManifestReset.WithLabelValues(user.Name).Set(float64(rateLimits.IntegrationManifest.Reset.Unix()))
	}

	log.Printf("Updated rate limits for user %s: Core=%d/%d, Search=%d/%d, GraphQL=%d/%d",
		user.Name,
		rateLimits.Core.Remaining, rateLimits.Core.Limit,
		rateLimits.Search.Remaining, rateLimits.Search.Limit,
		rateLimits.GraphQL.Remaining, rateLimits.GraphQL.Limit,
	)
}

// StartPolling starts a background goroutine that periodically updates rate limits
func (c *Collector) StartPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Do an initial update
	c.Update(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping rate limit polling")
			return
		case <-ticker.C:
			c.Update(ctx)
		}
	}
}
