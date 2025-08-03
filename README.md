# BJJ Store - Brazilian Jiu-Jitsu E-commerce Platform

A full-stack e-commerce application built with Vue.js frontend and Go backend, specializing in Brazilian Jiu-Jitsu gear and equipment.

## Live Demo

- **Frontend**: https://abundant-caring-production-f136.up.railway.app/
- **Admin Dashboard**: https://abundant-caring-production-f136.up.railway.app/admin
- **API Documentation**: https://your-backend.railway.app/swagger/index.html

## Admin Access

```
Email: admin@bjjstore.com
Password: admin123
```

## Test Payment Information

Use these test card numbers during checkout:

| Card Number | Scenario | CVV | Expiry |
|------------|----------|-----|--------|
| `4242 4242 4242 4242` | Successful Payment | 123 | 12/25 |
| `4000 0000 0000 0002` | Card Declined | 123 | 12/25 |
| `4000 0000 0000 0119` | Processing Error | 123 | 12/25 |

*Any name, future expiry date, and 3-4 digit CVV will work*

## Features

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


