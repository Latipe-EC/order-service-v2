## Latipe Order Service (Version 2)
Tech :
- Go (1.20)
- Fiber v2
- Gorm (MySQL v8)
- Redis
- RabbitMQ
- Go-Resty
- gRPC
- FiberPrometheus

Handling the purchase process in e-commerce using microservices architecture. The order service includes several function:
- CRUD orders data
- Statisticize commission, profit, order count,...

The order creation process involves two phases:
- Phase 1: Processes HTTP POST requests, retrieves data by making gRPC requests to other services, and sends messages (order_status:pending) into transaction service.
- Phase 2: Receives reply messages from transaction service and update order status (failed or success) into the database.


Server endpoints:
- API: http://localhost:5000/api/v2/orders
- Metrics (Prometheus): http://localhost:5000/metrics
- Health check: http://localhost:5000/readiness
- Fiber dashboard: http://localhost:5000/fiber/dashboard
- Status check: http://localhost:5000/health or http://localhost:5000/liveness
- Swagger: http://localhost:5000/swagger/
- gRPC: localhost:6000
- Basic Auth: admin:123123 (for metrics, health, status)
<hr>
<h4>Development by Tran Tien Dat</h4>
