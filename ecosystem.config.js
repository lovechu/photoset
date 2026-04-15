module.exports = {
  apps: [
    // 后端 Go 服务
    {
      name: 'photoset-backend',
      script: 'cmd/main.go',
      cwd: '.',
      interpreter: 'go',
      interpreter_args: 'run',
      env: {
        PHOTOSET_ENV: 'development',
      },
      watch: false, // Go 文件变化时重启
      ignore_watch: ['node_modules', 'frontend', 'frontend-admin'],
      max_memory_restart: '300M',
      error_file: 'logs/backend-error.log',
      out_file: 'logs/backend-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss',
      restart_delay: 1000,
      // 健康检查
      healthcheck: {
        url: 'http://localhost:8080/api/health',
        max_attempts: 3,
        timeout: 5000,
        unhealthy_uptime: 30000
      }
    },
    // 用户前端
    {
      name: 'photoset-frontend',
      script: 'npm',
      args: 'run dev',
      cwd: './frontend',
      watch: ['./src', 'package.json'],
      ignore_watch: ['node_modules'],
      env: {
        NODE_ENV: 'development',
        VITE_API_BASE: 'http://localhost:8080/api',
        PORT: 3000,
        HOST: '0.0.0.0'
      },
      max_memory_restart: '200M',
      error_file: '../logs/frontend-error.log',
      out_file: '../logs/frontend-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss',
      restart_delay: 1000
    },
    // 管理后台前端
    {
      name: 'photoset-admin',
      script: 'npm',
      args: 'run dev',
      cwd: './frontend-admin',
      watch: ['./src', 'package.json'],
      ignore_watch: ['node_modules'],
      env: {
        NODE_ENV: 'development',
        VITE_API_BASE: 'http://localhost:8080/api',
        PORT: 3001,
        HOST: '0.0.0.0'
      },
      max_memory_restart: '200M',
      error_file: '../logs/admin-error.log',
      out_file: '../logs/admin-out.log',
      log_date_format: 'YYYY-MM-DD HH:mm:ss',
      restart_delay: 1000
    }
  ]
};