# Auth0 Setup Guide for Hulta Foal Tracker

## 🔐 Initial Access and Configuration

### 1. Login to Auth0 Dashboard
- Go to [https://manage.auth0.com/](https://manage.auth0.com/)
- Use the provided credentials to log in

### 2. Application Configuration
#### Create/Select Application
- Navigate to "Applications" > "Applications"
- Click "Create Application" or select existing application
- Set Application Type: "Regular Web Application"
- Name: "Hulta Foal Tracker"

### 3. Application Settings
#### Allowed Callback URLs
```
https://hulta-foaltracker.app/callback
https://api.hulta-foaltracker.app/callback
http://localhost:3000/callback
```

#### Allowed Logout URLs
```
https://hulta-foaltracker.app
http://localhost:3000
```

#### Allowed Web Origins
```
https://hulta-foaltracker.app
http://localhost:3000
```

### 4. API Configuration
- Go to "APIs"
- Create a new API or select existing
- Name: "Hulta Foal Tracker API"
- Identifier: `https://api.hulta-foaltracker.app`

### 5. Permissions and Scopes
- Add API Permissions:
  * `read:horses`
  * `write:horses`
  * `read:pregnancy`
  * `write:pregnancy`

### 6. Roles Configuration
Create the following roles:
- `user`
- `admin`
- `farm_manager`

### 7. Connection Settings
- Enable Username/Password authentication
- Configure social logins if desired (Google, Apple, etc.)

### 8. Security Settings
- Enable Multi-Factor Authentication
- Configure Password Policy
- Set up Anomaly Detection

## 🛠 Environment Configuration
Update `.env` files with:
- Frontend `.env`:
  ```env
  VITE_AUTH0_DOMAIN=dev-r083cwkcv0pgz20x.eu.auth0.com
  VITE_AUTH0_CLIENT_ID=YOUR_CLIENT_ID
  VITE_AUTH0_AUDIENCE=https://api.hulta-foaltracker.app
  ```

- Backend `.env`:
  ```env
  AUTH0_AUDIENCE=https://api.hulta-foaltracker.app
  AUTH0_ALGORITHM=RS256
  AUTH0_DOMAIN=dev-r083cwkcv0pgz20x.eu.auth0.com
  AUTH0_ISSUER=https://dev-r083cwkcv0pgz20x.eu.auth0.com/
  ```

## 🚀 Next Steps
1. Verify all settings
2. Test login functionality
3. Configure additional security features

## 🆘 Troubleshooting
- Ensure all URLs are exactly matching
- Check that client ID and domain are correct
- Verify network and firewall settings

## 📝 Notes
- Keep client secrets and sensitive information confidential
- Regularly review and update security settings
