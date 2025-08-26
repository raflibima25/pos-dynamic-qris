# CLAUDE.md

## Sistem QRIS Dinamis untuk Kasir

### 🎯 Project Overview

Sistem QRIS dinamis yang memungkinkan kasir membuat QRIS dengan nominal yang sudah ditentukan, sehingga customer tidak perlu input manual nominal saat melakukan pembayaran. Sistem ini akan menyederhanakan flow transaksi dan meningkatkan efisiensi operasional toko.

### 💡 Core Problem Statement

Saat ini customer perlu input nominal secara manual ketika scan QRIS, yang dapat menyebabkan:

- Human error dalam input nominal
- Proses pembayaran yang lambat
- Antrian yang panjang di kasir
- Potensi kesalahan transaksi

### ✨ Solution

Sistem POS terintegrasi yang dapat:

- Generate QRIS dengan nominal yang sudah fixed per transaksi
- Customer langsung scan dan bayar tanpa input nominal
- Real-time monitoring status pembayaran
- Automated receipt generation

---

### 🏗️ Technical Architecture

#### Tech Stack

```
Frontend  : Next.js (React)
Backend   : Golang (Gin dan Gorm)
Database  : PostgreSQL
Payment   : Midtrans QRIS API
Hosting   : TBD (Docker ready)
```

#### System Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Next.js   │    │   Golang    │    │ PostgreSQL  │
│  Frontend   │───▶│   Backend   │───▶│  Database   │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                           ▼
                   ┌─────────────┐
                   │  Midtrans   │
                   │ Payment API │
                   └─────────────┘
```

---

### 🚀 Key Features

#### 🛒 Point of Sale (POS)

- **Product Management**: CRUD operations untuk master barang
- **Inventory Tracking**: Real-time stock monitoring
- **Cart Management**: Add/remove items, quantity adjustment
- **Discount System**: Per-item dan per-transaction discounts
- **Tax Calculation**: Automatic tax computation

#### 📱 Dynamic QRIS

- **Smart Generation**: Auto-generate QRIS dengan nominal fixed
- **Expiry Management**: Configurable timeout (default 10 minutes)
- **Auto Refresh**: Regenerate expired QRIS seamlessly
- **Multiple Format**: Support various QR code formats

#### 💳 Payment Processing

- **Midtrans Integration**: Native QRIS API integration
- **Real-time Status**: Live payment status monitoring
- **Webhook Handling**: Automated payment confirmation
- **Error Handling**: Robust payment failure management

#### 📊 Analytics & Reporting

- **Sales Dashboard**: Real-time sales analytics
- **Transaction History**: Complete audit trail
- **Financial Reports**: Daily/monthly/yearly reports
- **Export Options**: PDF, Excel, CSV formats

---

### 🔄 User Flow

#### Kasir Flow

```
1. Login → 2. Scan/Input Barang → 3. Review Total → 4. Generate QRIS
                ↓
8. Print Receipt ← 7. Payment Confirmed ← 6. Monitor Status ← 5. Show QRIS
```

#### Customer Flow

```
1. Scan QRIS → 2. Confirm Payment → 3. Payment Success
   (No manual input required!)
```

---

### 📋 Development Phases

#### Phase 1: Foundation (4-6 weeks)

**Core Infrastructure**

- [ ] Project setup (Golang + Next.js)
- [ ] Database schema design
- [ ] Authentication system
- [ ] Basic CRUD operations
- [ ] POS interface development

**Deliverables:**

- Working POS system
- Product management
- User authentication
- Basic transaction flow

#### Phase 2: Payment Integration (3-4 weeks)

**QRIS & Payment**

- [ ] Midtrans API integration
- [ ] QRIS generation logic
- [ ] Payment webhook handling
- [ ] Real-time status updates
- [ ] Error handling & retry logic

**Deliverables:**

- Dynamic QRIS generation
- Payment processing
- Status monitoring
- Transaction completion

#### Phase 3: Advanced Features (2-3 weeks)

**Enhancement & Optimization**

- [ ] Reporting dashboard
- [ ] Receipt generation
- [ ] Configuration management
- [ ] Performance optimization
- [ ] UI/UX improvements

**Deliverables:**

- Analytics dashboard
- Automated reporting
- System configuration
- Enhanced user experience

#### Phase 4: Testing & Deployment (2-3 weeks)

**Quality Assurance**

- [ ] Unit testing
- [ ] Integration testing
- [ ] Performance testing
- [ ] Security audit
- [ ] Production deployment

**Deliverables:**

- Tested application
- Documentation
- Production-ready deployment
- User training materials

---

### 🛠️ Development Guidelines

#### Code Standards

```
Backend (Golang):
- Use Go modules for dependency management
- Follow clean architecture principles
- Implement proper error handling
- Use struct validation
- Write comprehensive tests

Frontend (Next.js):
- Use TypeScript for type safety
- Implement responsive design
- Follow React best practices
- Use proper state management
- Optimize for performance
- Use Tailwind CSS for styling
```

#### API Design

```
REST API Endpoints:
GET    /api/products           # List products
POST   /api/transactions       # Create transaction
POST   /api/qris/generate      # Generate QRIS
GET    /api/qris/:id/status    # Check payment status
POST   /api/payments/callback  # Midtrans webhook
```

#### Database Design

```sql
Key Tables:
- users (authentication)
- products (inventory)
- transactions (sales records)
- transaction_items (cart items)
- payments (payment records)
- qris_codes (generated QR codes)
```

---

### 🔧 Setup Instructions

#### Prerequisites

```bash
# Required tools
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+
- Docker (optional)
```

#### Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=qris_pos
DB_USER=postgres
DB_PASS=your_password

# Midtrans
MIDTRANS_SERVER_KEY=your_server_key
MIDTRANS_CLIENT_KEY=your_client_key
MIDTRANS_ENVIRONMENT=sandbox

# App
JWT_SECRET=your_jwt_secret
APP_PORT=8080
```

#### Quick Start

```bash
# Backend setup
cd backend
go mod init qris-pos-backend
go mod tidy
go run main.go

# Frontend setup
cd frontend
npm install
npm run dev
```

---

### 🏗️ Project Structure & Architecture

#### Folder Structure

```
dynamic-qris/
├── backend/                    # Golang backend (Clean Architecture) (go1.24.0)
│   ├── cmd/
│   │   └── api/
│   │       └── main.go        # Application entry point
│   ├── internal/              # Private application code
│   │   ├── domain/            # Business logic layer
│   │   │   ├── entities/      # Business entities
│   │   │   ├── repositories/  # Repository interfaces
│   │   │   └── services/      # Business services
│   │   ├── infrastructure/    # External concerns
│   │   │   ├── database/      # Database implementations
│   │   │   ├── payment/       # Midtrans integration
│   │   │   ├── qrcode/        # QR code generation
│   │   │   └── config/        # Configuration management
│   │   ├── interfaces/        # Interface adapters
│   │   │   ├── http/          # HTTP handlers & routes
│   │   │   ├── middleware/    # HTTP middleware
│   │   │   └── dto/           # Data transfer objects
│   │   └── usecases/          # Application business rules
│   │       ├── product/       # Product use cases
│   │       ├── transaction/   # Transaction use cases
│   │       ├── payment/       # Payment use cases
│   │       └── auth/          # Authentication use cases
│   ├── pkg/                   # Public reusable packages
│   │   ├── logger/            # Logging utilities
│   │   ├── validator/         # Input validation
│   │   ├── errors/            # Custom error types
│   │   └── response/          # API response helpers
│   ├── migrations/            # Database migration files
│   ├── docs/                  # API documentation
│   ├── tests/                 # Test files
│   │   ├── integration/       # Integration tests
│   │   └── unit/              # Unit tests
│   ├── scripts/               # Build & deployment scripts
│   ├── .env.example           # Environment variables template
│   ├── Dockerfile             # Docker configuration
│   ├── go.mod                 # Go modules
│   └── Makefile               # Build commands
│
├── frontend/                  # Next.js frontend
│   ├── public/                # Static assets
│   ├── src/
│   │   ├── app/               # App router (Next.js latest)
│   │   │   ├── (auth)/        # Auth route group
│   │   │   ├── dashboard/     # Dashboard pages
│   │   │   ├── pos/           # POS interface
│   │   │   ├── products/      # Product management
│   │   │   ├── transactions/  # Transaction history
│   │   │   ├── reports/       # Analytics & reports
│   │   │   └── layout.tsx     # Root layout
│   │   ├── components/        # Reusable UI components
│   │   │   ├── ui/            # Base UI components
│   │   │   ├── forms/         # Form components
│   │   │   ├── layout/        # Layout components
│   │   │   └── charts/        # Chart components
│   │   ├── hooks/             # Custom React hooks
│   │   ├── lib/               # Utility libraries
│   │   │   ├── api.ts         # API client
│   │   │   ├── utils.ts       # Helper functions
│   │   │   └── validations.ts # Form validations
│   │   ├── store/             # State management (Zustand)
│   │   │   ├── auth.ts        # Auth store
│   │   │   ├── cart.ts        # Cart store
│   │   │   └── transaction.ts # Transaction store
│   │   ├── types/             # TypeScript definitions
│   │   └── constants/         # Application constants
│   ├── .env.local.example     # Environment variables
│   ├── next.config.js         # Next.js configuration
│   ├── tailwind.config.js     # Tailwind configuration
│   ├── tsconfig.json          # TypeScript configuration
│   └── package.json           # NPM dependencies
│
├── docker/                    # Docker configurations
│   ├── docker-compose.yml     # Development environment
│   └── docker-compose.prod.yml # Production environment
├── docs/                      # Project documentation
│   ├── api/                   # API documentation
│   ├── deployment/            # Deployment guides
│   └── user-guide/            # User manuals
├── scripts/                   # Utility scripts
│   ├── setup.sh              # Initial setup script
│   ├── deploy.sh             # Deployment script
│   └── backup.sh             # Database backup
├── .gitignore                 # Git ignore rules
├── README.md                  # Project overview
└── CLAUDE.md                  # This file
```

---

### ⚡ Development Commands

#### Backend Commands

```bash
# Development
make dev                 # Run development server with hot reload
make build              # Build production binary
make test               # Run all tests
make test-unit          # Run unit tests only
make test-integration   # Run integration tests
make test-coverage      # Run tests with coverage report

# Database
make db-migrate-up      # Run database migrations
make db-migrate-down    # Rollback database migrations
make db-seed            # Seed database with sample data
make db-reset           # Reset database (drop & recreate)

# Code Quality
make lint               # Run golangci-lint
make fmt                # Format code with gofmt
make vet                # Run go vet
make mod-tidy           # Clean up go.mod

# Docker
make docker-build       # Build Docker image
make docker-run         # Run with Docker Compose
make docker-down        # Stop Docker containers

# Utilities
make gen-docs           # Generate API documentation
make gen-mocks          # Generate test mocks
make clean              # Clean build artifacts
```

#### Frontend Commands

```bash
# Development
npm run dev             # Start development server
npm run build           # Build for production
npm run start           # Start production server
npm run lint            # Run ESLint
npm run lint:fix        # Fix ESLint issues

# Testing
npm run test            # Run Jest tests
npm run test:watch      # Run tests in watch mode
npm run test:coverage   # Run tests with coverage
npm run test:e2e        # Run Playwright E2E tests

# Type Checking
npm run type-check      # Run TypeScript compiler
npm run type-check:watch # Watch mode for type checking

# Code Quality
npm run format          # Format code with Prettier
npm run format:check    # Check code formatting
npm run analyze         # Bundle analyzer

# Storybook (if implemented)
npm run storybook       # Start Storybook server
npm run build-storybook # Build Storybook static files
```

#### Database Commands

```bash
# PostgreSQL Management
createdb qris_pos_dev          # Create development database
createdb qris_pos_test         # Create test database
dropdb qris_pos_dev           # Drop development database

# Migration Commands
migrate create -ext sql -dir migrations [name]  # Create new migration
migrate -path migrations -database "postgres://..." up    # Apply migrations
migrate -path migrations -database "postgres://..." down  # Rollback migrations
```

#### Docker Commands

```bash
# Development Environment
docker-compose up -d           # Start all services in background
docker-compose down           # Stop all services
docker-compose logs -f api    # Follow backend logs
docker-compose logs -f web    # Follow frontend logs

# Production Environment
docker-compose -f docker-compose.prod.yml up -d    # Start production
docker-compose -f docker-compose.prod.yml down     # Stop production

# Database Operations
docker-compose exec db psql -U postgres -d qris_pos  # Connect to database
docker-compose exec db pg_dump -U postgres qris_pos > backup.sql  # Backup
```

#### Testing Commands

```bash
# Backend Testing
go test ./...                          # Run all tests
go test ./internal/usecases/...        # Test specific package
go test -race ./...                    # Test with race detection
go test -bench=. ./...                 # Run benchmarks
go test -coverprofile=coverage.out ./... # Generate coverage

# Frontend Testing
npm run test -- --coverage            # Test with coverage
npm run test -- --watch              # Watch mode
npm run test -- ProductCard          # Test specific component
npm run test:e2e -- --headed         # E2E tests with browser

# Integration Testing
make test-integration                  # Backend integration tests
npm run test:e2e                     # Frontend E2E tests
```

---

### 🏛️ Clean Architecture Principles

#### Backend Architecture (Golang)

**Layer Separation:**

```
┌─────────────────────────────────────────────────────┐
│                   Frameworks & Drivers              │
│  (HTTP Handlers, Database, External APIs)          │
├─────────────────────────────────────────────────────┤
│                 Interface Adapters                  │
│        (Controllers, Presenters, Gateways)         │
├─────────────────────────────────────────────────────┤
│                 Application Business Rules           │
│                   (Use Cases)                      │
├─────────────────────────────────────────────────────┤
│               Enterprise Business Rules              │
│                   (Entities)                       │
└─────────────────────────────────────────────────────┘
```

**Dependency Rule:**

- Dependencies point inward
- Inner layers don't know about outer layers
- Use interfaces for dependency inversion

**Key Patterns:**

```go
// Repository Pattern
type ProductRepository interface {
    Create(ctx context.Context, product *domain.Product) error
    GetByID(ctx context.Context, id string) (*domain.Product, error)
    Update(ctx context.Context, product *domain.Product) error
    Delete(ctx context.Context, id string) error
}

// Use Case Pattern
type ProductUseCase struct {
    productRepo domain.ProductRepository
    logger      logger.Logger
}

// Entity Pattern
type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    CategoryID  string    `json:"category_id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Frontend Architecture (Next.js)

**Component Architecture:**

```typescript
// Feature-based organization
// components/pos/
├── PosLayout.tsx              # Layout component
├── ProductGrid.tsx            # Product selection
├── ShoppingCart.tsx           # Cart management
├── PaymentSummary.tsx         # Payment details
└── QRCodeDisplay.tsx          # QRIS display

// Custom Hooks Pattern
const useCart = () => {
  const [items, setItems] = useState<CartItem[]>([])

  const addItem = useCallback((product: Product) => {
    // Cart logic
  }, [])

  return { items, addItem, removeItem, total }
}

// Store Pattern (Zustand)
interface CartStore {
  items: CartItem[]
  total: number
  addItem: (product: Product) => void
  removeItem: (productId: string) => void
  clear: () => void
}
```

#### Best Practices Implementation

**Backend Best Practices:**

```go
// 1. Error Handling
type AppError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

// 2. Context Usage
func (u *ProductUseCase) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
    // Always pass context for cancellation and deadlines
    return u.productRepo.GetByID(ctx, id)
}

// 3. Validation
type CreateProductRequest struct {
    Name       string  `json:"name" validate:"required,min=1,max=100"`
    Price      float64 `json:"price" validate:"required,gt=0"`
    Stock      int     `json:"stock" validate:"required,gte=0"`
    CategoryID string  `json:"category_id" validate:"required,uuid"`
}

// 4. Configuration Management
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Midtrans MidtransConfig `mapstructure:"midtrans"`
}
```

**Frontend Best Practices:**

```typescript
// 1. Type Safety
interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
  categoryId: string;
}

// 2. Error Boundaries
class PosErrorBoundary extends Component<Props, State> {
  static getDerivedStateFromError(error: Error): State {
    return { hasError: true };
  }
}

// 3. Performance Optimization
const ProductCard = memo(({ product }: { product: Product }) => {
  return <div>{/* Product display */}</div>;
});

// 4. Custom Hooks for Business Logic
const usePayment = (transactionId: string) => {
  const [status, setStatus] = useState<PaymentStatus>("pending");

  useEffect(() => {
    // WebSocket connection for real-time updates
  }, [transactionId]);

  return { status, retry: retryPayment };
};
```

---

### 🔧 Development Workflow

#### Git Workflow

```bash
# Feature Development
git checkout -b feature/payment-integration
git add .
git commit -m "feat: implement Midtrans QRIS integration"
git push origin feature/payment-integration

# Code Review & Merge
# Create PR → Review → Merge to main
```

#### Code Review Checklist

- [ ] Follows clean architecture principles
- [ ] Has appropriate tests (unit + integration)
- [ ] Error handling implemented
- [ ] Input validation added
- [ ] Documentation updated
- [ ] No security vulnerabilities
- [ ] Performance considerations addressed

---

### 📊 Success Metrics

#### Performance KPIs

- **Transaction Speed**: < 30 seconds from scan to payment
- **QRIS Generation**: < 3 seconds
- **Payment Success Rate**: 99.9%
- **System Uptime**: 99%+
- **Concurrent Users**: 50+ simultaneous

#### Business KPIs

- **Customer Satisfaction**: Reduced checkout time by 50%
- **Error Reduction**: Eliminate manual input errors
- **Operational Efficiency**: Increase transaction throughput
- **Staff Productivity**: Streamlined workflow

---

### 🔒 Security Considerations

#### Data Protection

- **Encryption**: All sensitive data encrypted at rest and in transit
- **Authentication**: JWT-based secure authentication
- **Authorization**: Role-based access control
- **API Security**: Rate limiting and input validation
- **Payment Security**: PCI DSS compliant via Midtrans

#### Privacy Compliance

- **Data Minimization**: Only collect necessary data
- **Audit Trail**: Complete transaction logging
- **Access Control**: Restricted admin access
- **Data Retention**: Configurable data retention policy

---

### 🧪 Testing Strategy

#### Test Coverage

```
Unit Tests       : 80%+ coverage
Integration Tests: Payment flow, API endpoints
E2E Tests       : Complete user journeys
Performance     : Load testing with 100+ concurrent users
Security        : Vulnerability scanning
```

#### Testing Tools

- **Backend**: Go testing framework, testify
- **Frontend**: Jest, React Testing Library
- **E2E**: Playwright/Cypress
- **API**: Postman/Insomnia collections

---

### 📚 Documentation

#### Technical Docs

- [ ] API Documentation (OpenAPI/Swagger)
- [ ] Database Schema Documentation
- [ ] Deployment Guide
- [ ] Troubleshooting Guide

#### User Docs

- [ ] Admin User Manual
- [ ] Kasir Operation Guide
- [ ] Configuration Guide
- [ ] FAQ Document

---

### 🎁 Future Enhancements

#### Phase 2 Features

- **Multi-store Support**: Manage multiple outlets
- **Mobile App**: Dedicated mobile app for kasir
- **Advanced Analytics**: AI-powered insights
- **Loyalty Program**: Customer loyalty integration
- **Offline Mode**: Continue operations without internet

#### Integration Possibilities

- **Accounting Software**: QuickBooks, Jurnal integration
- **E-commerce**: Shopify, WooCommerce sync
- **WhatsApp Business**: Automated notifications
- **Printer Integration**: Direct receipt printing

---

### 👥 Team & Responsibilities

#### Development Team

```
Full-Stack Developer: Core development
UI/UX Designer     : User interface design
QA Engineer        : Testing & quality assurance
DevOps Engineer    : Infrastructure & deployment
```

#### Stakeholders

- **Business Owner**: Requirements & validation
- **Kasir/Staff**: User acceptance testing
- **IT Admin**: System maintenance
- **Finance**: Reporting requirements

---

### 📞 Support & Maintenance

#### Support Channels

- **Technical Issues**: GitHub Issues
- **Business Questions**: Direct communication
- **Emergency**: Phone support during business hours

#### Maintenance Schedule

- **Daily**: Automated backups
- **Weekly**: Performance monitoring
- **Monthly**: Security updates
- **Quarterly**: Feature updates

---

### 📄 License & Legal

#### Software License

- **Open Source Components**: MIT/Apache licenses
- **Custom Code**: Proprietary license
- **Third-party APIs**: Respective vendor licenses

#### Compliance

- **PCI DSS**: Payment card industry standards
- **Bank Indonesia**: QRIS compliance
- **Local Regulations**: Indonesian fintech regulations

---

### 📝 Changelog

#### Version History

```
v1.0.0 (Planned)
- Initial release
- Core POS functionality
- QRIS integration
- Basic reporting

v1.1.0 (Future)
- Multi-store support
- Advanced analytics
- Mobile optimization
```

---

### 🤝 Contributing

#### How to Contribute

1. Fork the repository
2. Create feature branch
3. Make your changes
4. Add tests
5. Submit pull request

#### Code Review Process

- Automated testing must pass
- Code review by senior developer
- Security review for payment-related changes
- Documentation updates

---

_Last Updated: August 26, 2025_
_Document Version: 1.0_
