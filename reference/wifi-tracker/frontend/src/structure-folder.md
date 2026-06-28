```
src/
├── components/
│ ├── ui/ # Reusable UI components
│ │ ├── Button.tsx
│ │ ├── Input.tsx
│ │ ├── Modal.tsx
│ │ ├── DataTable.tsx
│ │ └── index.ts
│ ├── forms/ # Form components
│ │ ├── WifiRegistrationForm.tsx
│ │ ├── UserForm.tsx
│ │ └── index.ts
│ ├── layout/ # Layout components
│ │ ├── SideBar.tsx
│ │ ├── Header.tsx
│ │ └── index.ts
│ └── features/ # Feature-specific components
│ ├── wifi/
│ │ ├── WifiList.tsx
│ │ ├── WifiCard.tsx
│ │ └── WifiFilters.tsx
│ └── users/
│ ├── UserList.tsx
│ ├── UserCard.tsx
│ └── UserFilters.tsx
├── hooks/ # Custom React hooks
│ ├── useWifiRegistration.ts
│ ├── useUsers.ts
│ ├── useAuth.ts
│ └── index.ts
├── services/ # API services
│ ├── api.ts
│ ├── wifiService.ts
│ ├── userService.ts
│ ├── authService.ts
│ └── index.ts
├── types/ # TypeScript type definitions
│ ├── wifi.ts
│ ├── user.ts
│ ├── api.ts
│ └── index.ts
├── utils/ # Utility functions
│ ├── auth.ts
│ ├── validation.ts
│ ├── constants.ts
│ └── index.ts
├── pages/ # Page components
│ ├── WifiRegistrationPage.tsx
│ ├── DashboardPage.tsx
│ ├── UsersPage.tsx
│ └── NotFoundPage.tsx
└── routes/ # Route definitions
├── index.tsx
├── wifi.tsx
└── users.tsx
```
