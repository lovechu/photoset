module.exports = {
  apps: [
    {
      name: 'photoset-backend',
      script: './main',
      cwd: '/mnt/c/Users/ichuy/WorkBuddy/20260408115223/backend',
      env: {
        NODE_ENV: 'development'
      },
      watch: false,
      autorestart: true
    },
    {
      name: 'photoset-frontend',
      script: 'npm',
      args: 'run dev',
      cwd: '/mnt/c/Users/ichuy/WorkBuddy/20260408115223/backend/frontend',
      interpreter: 'bash',
      watch: false,
      autorestart: true
    },
    {
      name: 'photoset-admin',
      script: 'npm',
      args: 'run dev',
      cwd: '/mnt/c/Users/ichuy/WorkBuddy/20260408115223/backend/frontend-admin',
      interpreter: 'bash',
      watch: false,
      autorestart: true
    }
  ]
}
