# 🥋 BJJ Store - Brazilian Jiu-Jitsu E-commerce Platform

A full-stack e-commerce application built with Vue.js frontend and Go backend, specializing in Brazilian Jiu-Jitsu gear and equipment.

## 🌐 Live Demo

- **Frontend**: https://your-frontend.railway.app
- **Admin Dashboard**: https://your-frontend.railway.app/admin
- **API Documentation**: https://your-backend.railway.app/swagger/index.html

## 👤 Admin Access

```
Email: admin@bjjstore.com
Password: admin123
```

## 💳 Test Payment Information

Use these test card numbers during checkout:

| Card Number | Scenario | CVV | Expiry |
|------------|----------|-----|--------|
| `4242 4242 4242 4242` | ✅ Successful Payment | 123 | 12/25 |
| `4000 0000 0000 0002` | ❌ Card Declined | 123 | 12/25 |
| `4000 0000 0000 0119` | ⚠️ Processing Error | 123 | 12/25 |

*Any name, future expiry date, and 3-4 digit CVV will work*

## 🚀 Features

### Customer Features
- **Product Catalog**: Browse BJJ gis, belts, and equipment
- **Product Search & Filtering**: Find products quickly
- **Product Details**: View images, descriptions, sizes, and pricing
- **Shopping Cart**: Add/remove items with size and quantity selection
- **Guest Checkout**: Complete purchases without account creation
- **Payment Processing**: Secure mock payment gateway
- **Order Tracking**: Track order status with order number
- **Email Order History**: Retrieve orders by email address
- **Mobile Responsive**: Optimized for all device sizes

### Admin Features
- **Admin Dashboard**: Real-time statistics and recent orders overview
- **Product Management**: Add, edit, delete products with image uploads
- **Order Management**: View all orders, update order status
- **Image Upload**: MinIO-powered cloud storage for product images
- **Auto-refresh**: Dashboard updates every 30 seconds
- **Role-based Access**: Different permission levels for admin users
- **JWT Authentication**: Secure admin authentication system

## 🛠 Tech Stack

### Frontend
- **Vue.js 3** - Progressive JavaScript framework
- **TypeScript** - Type-safe JavaScript
- **Pinia** - State management
- **Vue Router** - Client-side routing
- **Tailwind CSS** - Utility-first CSS framework
- **Axios** - HTTP client
- **Vite** - Build tool and dev server

### Backend
- **Go** - Systems programming language
- **Gin** - HTTP web framework
- **GORM** - Object-relational mapping
- **PostgreSQL** - Relational database
- **JWT** - JSON Web Token authentication
- **MinIO** - S3-compatible object storage
- **Swagger** - API documentation

### Infrastructure
- **Railway** - Cloud deployment platform
- **Docker** - Containerization
- **GitHub Actions** - CI/CD pipeline
- **Nginx** - Web server and reverse proxy

## 📋 API Endpoints

### Public Endpoints
```
GET    /api/health              # Health check
GET    /api/products            # List all products
GET    /api/products/:id        # Get product details
POST   /api/orders              # Create new order
GET    /api/orders/track/:orderNumber # Track order
POST   /api/payment/process     # Process payment
```

### Admin Endpoints
```
POST   /api/admin/auth/login    # Admin login
GET    /api/admin/profile       # Admin profile
GET    /api/admin/stats         # Dashboard statistics
POST   /api/admin/products      # Create product
PUT    /api/admin/products/:id  # Update product
DELETE /api/admin/products/:id  # Delete product
GET    /api/admin/orders        # List all orders
PUT    /api/admin/orders/:id/status # Update order status
POST   /api/admin/upload/image  # Upload product image
```

## 🏗 Architecture

```
┌─────────────────┐    HTTP/HTTPS    ┌─────────────────┐
│                 │ ◄──────────────► │                 │
│  Vue.js Client  │                  │   Go Backend    │
│                 │                  │                 │
└─────────────────┘                  └─────────────────┘
                                              │
                                              ▼
                                     ┌─────────────────┐
                                     │   PostgreSQL    │
                                     │    Database     │
                                     └─────────────────┘
                                              │
                                              ▼
                                     ┌─────────────────┐
                                     │  MinIO Storage  │
                                     │  (Product Images)│
                                     └─────────────────┘
```

## 🚀 Local Development

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Docker (optional)

### Backend Setup
```bash
cd backend
go mod download
cp .env.example .env
# Edit .env with your database credentials
go run main.go
```

### Frontend Setup
```bash
cd frontend
npm install
npm run dev
```

### Database Setup
```bash
# Using Docker
docker-compose up -d postgres

# Or install PostgreSQL locally and create database
createdb bjj_store
```

## 📦 Deployment

The application is deployed on Railway with the following services:

1. **PostgreSQL Database** - Managed database service
2. **Go Backend** - API server with health checks
3. **Vue.js Frontend** - Static site with Nginx
4. **MinIO Storage** - Object storage for images

### Environment Variables

#### Backend
```
DATABASE_HOST=<database-host>
DATABASE_PORT=5432
DATABASE_USER=<database-user>
DATABASE_PASSWORD=<database-password>
DATABASE_NAME=<database-name>
DATABASE_SSL_MODE=require
SERVER_PORT=8080
JWT_SECRET=<your-jwt-secret>
MINIO_ENDPOINT=<minio-endpoint>
```

#### Frontend
```
VITE_API_BASE_URL=<backend-api-url>
```

## 🧪 Testing

### Payment Testing
Use the test card numbers provided above to simulate different payment scenarios.

### Admin Testing
1. Login to admin dashboard with provided credentials
2. Add test products with images
3. View and manage customer orders
4. Test order status updates

### Customer Flow Testing
1. Browse product catalog
2. Add items to cart with different sizes
3. Complete checkout process
4. Track order with order number

## 📱 Mobile Support

The application is fully responsive and optimized for:
- Desktop browsers (Chrome, Firefox, Safari, Edge)
- Tablet devices (iPad, Android tablets)
- Mobile phones (iOS Safari, Android Chrome)

## 🔒 Security Features

- **JWT Authentication** for admin access
- **Input validation** on all forms
- **SQL injection protection** with GORM
- **XSS prevention** with proper escaping
- **CORS configuration** for API security
- **Environment variable management** for secrets

## 👥 User Roles

### Customer (Public)
- Browse products
- Add to cart
- Checkout and pay
- Track orders

### Admin
- All customer features plus:
- Product management
- Order management
- Image uploads
- Dashboard analytics

## 📞 Support

For questions or issues:
- Check the API documentation at `/swagger/index.html`
- Review the admin dashboard for order management
- Use test payment cards for development

## 🎯 Project Goals

This project demonstrates:
- Full-stack web development skills
- Modern JavaScript/TypeScript frontend development
- Backend API development with Go
- Database design and management
- Cloud deployment and DevOps
- E-commerce functionality implementation
- Payment processing integration
- Admin dashboard development
- Responsive web design
- Security best practices

---

**Built with ❤️ for Brazilian Jiu-Jitsu enthusiasts**