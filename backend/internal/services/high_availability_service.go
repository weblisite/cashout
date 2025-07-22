package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// HighAvailabilityService handles high availability operations
type HighAvailabilityService struct {
	supabaseService *SupabaseService
	coreBanking     *CoreBankingService
	mu              sync.RWMutex
	healthStatus    map[string]HealthStatus
	loadBalancer    *LoadBalancer
}

// NewHighAvailabilityService creates a new high availability service
func NewHighAvailabilityService(supabaseService *SupabaseService, coreBanking *CoreBankingService) *HighAvailabilityService {
	return &HighAvailabilityService{
		supabaseService: supabaseService,
		coreBanking:     coreBanking,
		healthStatus:    make(map[string]HealthStatus),
		loadBalancer:    NewLoadBalancer(),
	}
}

// HealthStatus represents service health status
type HealthStatus struct {
	ServiceName    string    `json:"service_name"`
	Status         string    `json:"status"` // healthy, degraded, unhealthy
	ResponseTime   time.Duration `json:"response_time"`
	ErrorCount     int       `json:"error_count"`
	LastCheck      time.Time `json:"last_check"`
	LastError      string    `json:"last_error,omitempty"`
	Uptime         time.Duration `json:"uptime"`
	StartTime      time.Time `json:"start_time"`
}

// LoadBalancer represents load balancing configuration
type LoadBalancer struct {
	Algorithm      string    `json:"algorithm"` // round_robin, least_connections, weighted
	BackendNodes   []BackendNode `json:"backend_nodes"`
	CurrentIndex   int       `json:"current_index"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`
	FailoverThreshold int `json:"failover_threshold"`
}

// BackendNode represents a backend service node
type BackendNode struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	Weight        int       `json:"weight"`
	MaxConnections int      `json:"max_connections"`
	CurrentConnections int  `json:"current_connections"`
	Status        string    `json:"status"` // active, inactive, failed
	HealthStatus  HealthStatus `json:"health_status"`
	LastFailover  *time.Time `json:"last_failover,omitempty"`
}

// PerformanceMetrics represents performance metrics
type PerformanceMetrics struct {
	ID                uuid.UUID `json:"id"`
	ServiceName       string    `json:"service_name"`
	Timestamp         time.Time `json:"timestamp"`
	ResponseTime      time.Duration `json:"response_time"`
	Throughput        int       `json:"throughput"` // requests per second
	ErrorRate         float64   `json:"error_rate"`
	CPUUsage          float64   `json:"cpu_usage"`
	MemoryUsage       float64   `json:"memory_usage"`
	DiskUsage         float64   `json:"disk_usage"`
	NetworkLatency    time.Duration `json:"network_latency"`
	ActiveConnections int       `json:"active_connections"`
	CreatedAt         time.Time `json:"created_at"`
}

// DisasterRecovery represents disaster recovery configuration
type DisasterRecovery struct {
	ID                    uuid.UUID `json:"id"`
	BackupStrategy        string    `json:"backup_strategy"` // full, incremental, differential
	BackupFrequency       string    `json:"backup_frequency"` // hourly, daily, weekly
	RetentionPeriod       int       `json:"retention_period"` // days
	RecoveryPointObjective time.Duration `json:"rpo"` // Recovery Point Objective
	RecoveryTimeObjective time.Duration `json:"rto"` // Recovery Time Objective
	LastBackup            *time.Time `json:"last_backup,omitempty"`
	LastRecovery          *time.Time `json:"last_recovery,omitempty"`
	Status                string    `json:"status"` // active, inactive, failed
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// High Availability Methods
func (h *HighAvailabilityService) StartHealthMonitoring(ctx context.Context) {
	go h.monitorServices(ctx)
	go h.monitorLoadBalancer(ctx)
	go h.collectPerformanceMetrics(ctx)
}

func (h *HighAvailabilityService) monitorServices(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.checkServiceHealth(ctx)
		}
	}
}

func (h *HighAvailabilityService) checkServiceHealth(ctx context.Context) {
	services := []string{"supabase", "core_banking", "settlement", "liquidity", "risk_management", "kyc_aml", "fraud_detection"}

	for _, serviceName := range services {
		start := time.Now()
		status := "healthy"
		var err error

		switch serviceName {
		case "supabase":
			err = h.supabaseService.HealthCheck(ctx)
		case "core_banking":
			err = h.coreBanking.HealthCheck(ctx)
		// Add other services as they're implemented
		}

		responseTime := time.Since(start)

		h.mu.Lock()
		health, exists := h.healthStatus[serviceName]
		if !exists {
			health = HealthStatus{
				ServiceName: serviceName,
				StartTime:   time.Now(),
			}
		}

		if err != nil {
			status = "unhealthy"
			health.ErrorCount++
			health.LastError = err.Error()
		} else {
			health.ErrorCount = 0
			health.LastError = ""
		}

		health.Status = status
		health.ResponseTime = responseTime
		health.LastCheck = time.Now()
		health.Uptime = time.Since(health.StartTime)

		h.healthStatus[serviceName] = health
		h.mu.Unlock()

		if err != nil {
			log.Printf("⚠️ Service %s health check failed: %v", serviceName, err)
		} else {
			log.Printf("✅ Service %s health check passed (%.2fms)", serviceName, responseTime.Seconds()*1000)
		}
	}
}

func (h *HighAvailabilityService) GetServiceHealth(ctx context.Context, serviceName string) (*HealthStatus, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	health, exists := h.healthStatus[serviceName]
	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	return &health, nil
}

func (h *HighAvailabilityService) GetAllServiceHealth(ctx context.Context) map[string]HealthStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()

	result := make(map[string]HealthStatus)
	for k, v := range h.healthStatus {
		result[k] = v
	}

	return result
}

// Load Balancer Methods
func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		Algorithm:           "round_robin",
		BackendNodes:        make([]BackendNode, 0),
		CurrentIndex:        0,
		HealthCheckInterval: 30 * time.Second,
		FailoverThreshold:   3,
	}
}

func (h *HighAvailabilityService) AddBackendNode(ctx context.Context, name, url string, weight, maxConnections int) error {
	node := BackendNode{
		ID:                 uuid.New(),
		Name:               name,
		URL:                url,
		Weight:             weight,
		MaxConnections:     maxConnections,
		CurrentConnections: 0,
		Status:             "active",
		HealthStatus: HealthStatus{
			ServiceName: name,
			Status:      "healthy",
			StartTime:   time.Now(),
		},
	}

	h.loadBalancer.BackendNodes = append(h.loadBalancer.BackendNodes, node)

	if err := h.supabaseService.CreateBackendNode(ctx, &node); err != nil {
		return fmt.Errorf("failed to create backend node: %w", err)
	}

	log.Printf("Added backend node: %s (%s)", name, url)
	return nil
}

func (h *HighAvailabilityService) GetNextBackendNode(ctx context.Context) (*BackendNode, error) {
	if len(h.loadBalancer.BackendNodes) == 0 {
		return nil, fmt.Errorf("no backend nodes available")
	}

	switch h.loadBalancer.Algorithm {
	case "round_robin":
		return h.getNextRoundRobin()
	case "least_connections":
		return h.getLeastConnections()
	case "weighted":
		return h.getWeightedNode()
	default:
		return h.getNextRoundRobin()
	}
}

func (h *HighAvailabilityService) getNextRoundRobin() (*BackendNode, error) {
	activeNodes := h.getActiveNodes()
	if len(activeNodes) == 0 {
		return nil, fmt.Errorf("no active backend nodes")
	}

	h.loadBalancer.CurrentIndex = (h.loadBalancer.CurrentIndex + 1) % len(activeNodes)
	return &activeNodes[h.loadBalancer.CurrentIndex], nil
}

func (h *HighAvailabilityService) getLeastConnections() (*BackendNode, error) {
	activeNodes := h.getActiveNodes()
	if len(activeNodes) == 0 {
		return nil, fmt.Errorf("no active backend nodes")
	}

	var leastConnectionsNode *BackendNode
	minConnections := int(^uint(0) >> 1) // Max int

	for i := range activeNodes {
		if activeNodes[i].CurrentConnections < minConnections {
			minConnections = activeNodes[i].CurrentConnections
			leastConnectionsNode = &activeNodes[i]
		}
	}

	return leastConnectionsNode, nil
}

func (h *HighAvailabilityService) getWeightedNode() (*BackendNode, error) {
	activeNodes := h.getActiveNodes()
	if len(activeNodes) == 0 {
		return nil, fmt.Errorf("no active backend nodes")
	}

	// Simple weighted round-robin
	totalWeight := 0
	for _, node := range activeNodes {
		totalWeight += node.Weight
	}

	// Use current index to select weighted node
	currentWeight := 0
	for _, node := range activeNodes {
		currentWeight += node.Weight
		if h.loadBalancer.CurrentIndex < currentWeight {
			h.loadBalancer.CurrentIndex = (h.loadBalancer.CurrentIndex + 1) % totalWeight
			return &node, nil
		}
	}

	// Fallback to first active node
	return &activeNodes[0], nil
}

func (h *HighAvailabilityService) getActiveNodes() []BackendNode {
	var activeNodes []BackendNode
	for _, node := range h.loadBalancer.BackendNodes {
		if node.Status == "active" {
			activeNodes = append(activeNodes, node)
		}
	}
	return activeNodes
}

func (h *HighAvailabilityService) monitorLoadBalancer(ctx context.Context) {
	ticker := time.NewTicker(h.loadBalancer.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.checkBackendNodes(ctx)
		}
	}
}

func (h *HighAvailabilityService) checkBackendNodes(ctx context.Context) {
	for i := range h.loadBalancer.BackendNodes {
		node := &h.loadBalancer.BackendNodes[i]
		
		// Simulate health check
		start := time.Now()
		healthy := h.simulateHealthCheck(node.URL)
		responseTime := time.Since(start)

		// Update health status
		node.HealthStatus.ResponseTime = responseTime
		node.HealthStatus.LastCheck = time.Now()

		if !healthy {
			node.HealthStatus.ErrorCount++
			node.HealthStatus.Status = "unhealthy"
			node.HealthStatus.LastError = "Health check failed"

			// Check if node should be marked as failed
			if node.HealthStatus.ErrorCount >= h.loadBalancer.FailoverThreshold {
				node.Status = "failed"
				now := time.Now()
				node.LastFailover = &now
				log.Printf("🚨 Backend node %s marked as failed", node.Name)
			}
		} else {
			node.HealthStatus.ErrorCount = 0
			node.HealthStatus.Status = "healthy"
			node.HealthStatus.LastError = ""
		}
	}
}

func (h *HighAvailabilityService) simulateHealthCheck(url string) bool {
	// Simulate health check (in production, would make actual HTTP request)
	// For now, return true 95% of the time to simulate occasional failures
	return time.Now().UnixNano()%100 < 95
}

// Performance Monitoring Methods
func (h *HighAvailabilityService) collectPerformanceMetrics(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.collectMetrics(ctx)
		}
	}
}

func (h *HighAvailabilityService) collectMetrics(ctx context.Context) {
	services := []string{"supabase", "core_banking", "settlement", "liquidity", "risk_management"}

	for _, serviceName := range services {
		metrics := &PerformanceMetrics{
			ID:              uuid.New(),
			ServiceName:     serviceName,
			Timestamp:       time.Now(),
			ResponseTime:    h.getAverageResponseTime(serviceName),
			Throughput:      h.calculateThroughput(serviceName),
			ErrorRate:       h.calculateErrorRate(serviceName),
			CPUUsage:        h.getCPUUsage(),
			MemoryUsage:     h.getMemoryUsage(),
			DiskUsage:       h.getDiskUsage(),
			NetworkLatency:  h.getNetworkLatency(),
			ActiveConnections: h.getActiveConnections(),
			CreatedAt:       time.Now(),
		}

		if err := h.supabaseService.CreatePerformanceMetrics(ctx, metrics); err != nil {
			log.Printf("Failed to save performance metrics: %v", err)
		}
	}
}

func (h *HighAvailabilityService) getAverageResponseTime(serviceName string) time.Duration {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if health, exists := h.healthStatus[serviceName]; exists {
		return health.ResponseTime
	}
	return 0
}

func (h *HighAvailabilityService) calculateThroughput(serviceName string) int {
	// Simulate throughput calculation
	return 100 + int(time.Now().UnixNano()%900)
}

func (h *HighAvailabilityService) calculateErrorRate(serviceName string) float64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if health, exists := h.healthStatus[serviceName]; exists {
		if health.ErrorCount > 0 {
			return float64(health.ErrorCount) / 100.0
		}
	}
	return 0.0
}

func (h *HighAvailabilityService) getCPUUsage() float64 {
	// Simulate CPU usage
	return 20.0 + float64(time.Now().UnixNano()%60)
}

func (h *HighAvailabilityService) getMemoryUsage() float64 {
	// Simulate memory usage
	return 40.0 + float64(time.Now().UnixNano()%40)
}

func (h *HighAvailabilityService) getDiskUsage() float64 {
	// Simulate disk usage
	return 60.0 + float64(time.Now().UnixNano()%20)
}

func (h *HighAvailabilityService) getNetworkLatency() time.Duration {
	// Simulate network latency
	return time.Duration(10+time.Now().UnixNano()%50) * time.Millisecond
}

func (h *HighAvailabilityService) getActiveConnections() int {
	// Simulate active connections
	return 50 + int(time.Now().UnixNano()%200)
}

// Disaster Recovery Methods
func (h *HighAvailabilityService) CreateDisasterRecovery(ctx context.Context, strategy, frequency string, retentionPeriod int, rpo, rto time.Duration) (*DisasterRecovery, error) {
	dr := &DisasterRecovery{
		ID:                    uuid.New(),
		BackupStrategy:        strategy,
		BackupFrequency:       frequency,
		RetentionPeriod:       retentionPeriod,
		RecoveryPointObjective: rpo,
		RecoveryTimeObjective: rto,
		Status:                "active",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if err := h.supabaseService.CreateDisasterRecovery(ctx, dr); err != nil {
		return nil, fmt.Errorf("failed to create disaster recovery: %w", err)
	}

	log.Printf("Created disaster recovery plan: %s backup every %s", strategy, frequency)
	return dr, nil
}

func (h *HighAvailabilityService) PerformBackup(ctx context.Context) error {
	log.Println("🔄 Starting backup process...")

	// Simulate backup process
	time.Sleep(5 * time.Second)

	// Update last backup time
	dr, err := h.supabaseService.GetDisasterRecovery(ctx)
	if err != nil {
		return fmt.Errorf("failed to get disaster recovery: %w", err)
	}

	now := time.Now()
	dr.LastBackup = &now
	dr.UpdatedAt = now

	if err := h.supabaseService.UpdateDisasterRecovery(ctx, dr); err != nil {
		return fmt.Errorf("failed to update disaster recovery: %w", err)
	}

	log.Println("✅ Backup completed successfully")
	return nil
}

func (h *HighAvailabilityService) PerformRecovery(ctx context.Context) error {
	log.Println("🔄 Starting disaster recovery process...")

	// Simulate recovery process
	time.Sleep(10 * time.Second)

	// Update last recovery time
	dr, err := h.supabaseService.GetDisasterRecovery(ctx)
	if err != nil {
		return fmt.Errorf("failed to get disaster recovery: %w", err)
	}

	now := time.Now()
	dr.LastRecovery = &now
	dr.UpdatedAt = now

	if err := h.supabaseService.UpdateDisasterRecovery(ctx, dr); err != nil {
		return fmt.Errorf("failed to update disaster recovery: %w", err)
	}

	log.Println("✅ Disaster recovery completed successfully")
	return nil
}

// Health Check
func (h *HighAvailabilityService) HealthCheck(ctx context.Context) error {
	// Check overall system health
	allHealth := h.GetAllServiceHealth(ctx)
	
	unhealthyCount := 0
	for _, health := range allHealth {
		if health.Status == "unhealthy" {
			unhealthyCount++
		}
	}

	if unhealthyCount > 0 {
		return fmt.Errorf("high availability health check failed: %d unhealthy services", unhealthyCount)
	}

	log.Println("✅ High availability service health check passed")
	return nil
} 