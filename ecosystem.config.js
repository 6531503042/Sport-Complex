module.exports = {
    apps: [
      {
        name: "auth-service",
        script: "go",
        args: "run main.go ./env/dev/.env.auth",
        cwd: "/Users/bricksocoolxd/GitHub-Initial/Sport-Complex/backend", // Replace with your actual backend path
        env: {
          NODE_ENV: "development",
        },
      },
      
      {
        name: "booking-service",
        script: "go",
        args: "run main.go ./env/dev/.env.booking",
        cwd: "/Users/bricksocoolxd/GitHub-Initial/Sport-Complex/backend", // Replace with your actual backend path
        env: {
          NODE_ENV: "development",
        },
      },
      {
        name: "user-service",
        script: "go",
        args: "run main.go ./env/dev/.env.user",
        cwd: "/Users/bricksocoolxd/GitHub-Initial/Sport-Complex/backend", // Replace with your actual backend path
        env: {
          NODE_ENV: "development",
        },
      },
      {
        name: "facility-service",
        script: "go",
        args: "run main.go ./env/dev/.env.facility",
        cwd: "/Users/bricksocoolxd/GitHub-Initial/Sport-Complex/backend", // Replace with your actual backend path
        env: {
          NODE_ENV: "development",
        },
      },
      {
        name: "payment-service",
        script: "go",
        args: "run main.go ./env/dev/.env.payment",
        cwd: "/Users/bricksocoolxd/GitHub-Initial/Sport-Complex/backend", // Replace with your actual backend path
        env: {
          NODE_ENV: "development",
        },
      },
    ],
  };
  