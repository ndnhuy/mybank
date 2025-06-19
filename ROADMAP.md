## Realistic Scaling Evolution Roadmap

### Phase 1: Add Caching Layer (Weeks 1-2)
**Business Problem:** 
- Response times getting slower (users complaining)
- Database CPU hitting 80%+ during peak hours
- Same account info being queried repeatedly

**Technical Problem:**
- Every API call hits database
- Simple account balance checks causing expensive DB queries

**Implementation:**
- Add Redis for account balance caching
- Cache account info for 5-10 minutes
- Cache-aside pattern implementation

**New Distributed Systems Challenges:**
- **Cache invalidation**: When to clear account balance after transfer?
- **Cache consistency**: What if cache shows old balance?
- **Cache warming**: How to handle cache misses during startup?

**Metrics to Track:**
- Cache hit ratio
- Database query reduction
- Response time improvement

---

### Phase 2: Horizontal App Scaling (Weeks 3-4)
**Business Problem:**
- App server CPU hitting 100% during peak hours
- Limited concurrent users (few hundred max)
- Single point of failure for the entire system

**Technical Problem:**
- Single JVM can only handle limited concurrent requests
- All traffic funneled through one instance

**Implementation:**
- Deploy 2-3 instances of your Java app
- Add load balancer (nginx or cloud LB)
- Make your app stateless

**New Distributed Systems Challenges:**
- **Load balancing**: How to distribute traffic fairly?
- **Session affinity**: Do you need sticky sessions?
- **Shared cache access**: Multiple apps hitting same Redis
- **Health checks**: How to detect when an app instance is down?

**Metrics to Track:**
- Requests per second across instances
- CPU utilization per instance
- Load balancer health checks

---

### Phase 3: Distributed Caching (Weeks 5-7)
**Business Problem:**
- Redis becoming memory bottleneck (8GB limit hit)
- Cache becomes single point of failure
- Hot accounts (celebrities/businesses) overloading single cache

**Technical Problem:**
- Single Redis instance memory limits
- Cache availability issues

**Implementation:**
- Redis cluster with 3-5 nodes
- Consistent hashing for key distribution
- Cache partitioning by account ID

**New Distributed Systems Challenges:**
- **Hot key problem**: Popular accounts overloading one cache node
- **Cache partitioning**: How to distribute keys evenly?
- **Cross-cache invalidation**: Account A001 and A002 on different cache nodes
- **Cache cluster failover**: What happens when cache node dies?

**Metrics to Track:**
- Memory usage per cache node
- Hot key detection
- Cache cluster health

---

### Phase 4: Database Read Replicas (Weeks 8-10)
**Business Problem:**
- Database becoming bottleneck even with caching
- Need for 99.9% availability (planned maintenance downtime unacceptable)
- Read queries still hitting primary database hard

**Technical Problem:**
- Single database handling both reads and writes
- Database maintenance requires downtime

**Implementation:**
- 2-3 MySQL read replicas
- Route read queries to replicas
- Write queries still go to primary

**New Distributed Systems Challenges:**
- **Read replica lag**: Balance just transferred but read shows old balance
- **Read-write splitting**: Which queries go where?
- **Replica failover**: What if read replica goes down?
- **Eventual consistency**: How to handle lag between primary and replicas?

**Real-World Example Problem:**
```
User transfers $100 at 10:00:01
User checks balance at 10:00:02 (hits replica with 3-second lag)
User sees old balance and thinks transfer failed!
```

**Metrics to Track:**
- Replication lag
- Read/write query distribution
- Replica availability

---

### Phase 5: Database Sharding (Weeks 11-14)
**Business Problem:**
- Database write performance hitting limits (1000+ TPS)
- Database storage growing beyond single machine capacity (1TB+)
- Geographic expansion requiring data locality

**Technical Problem:**
- Single primary database can't handle write load
- Cross-region latency for global users

**Implementation:**
- Shard accounts across 2-4 database instances
- Consistent hashing by account ID
- Update your BankService to route to correct shard

**New Distributed Systems Challenges:**
- **Cross-shard transactions**: Transfer between accounts on different shards
- **Shard rebalancing**: What happens when you add new shard?
- **Global queries**: "Show all accounts" becomes complex
- **Hot shard problem**: Some shards busier than others

**This is where you'd finally need the 2PC or Saga patterns we discussed earlier!**

**Metrics to Track:**
- Transactions per shard
- Cross-shard transaction percentage
- Shard load distribution

---

## The Natural Learning Progression

Each phase teaches you:
1. **Caching**: Data consistency fundamentals
2. **App scaling**: Load distribution and stateless design
3. **Distributed caching**: Partitioning and hot key problems
4. **Read replicas**: Eventual consistency and CAP theorem
5. **Sharding**: Distributed transactions and the need for Saga/2PC