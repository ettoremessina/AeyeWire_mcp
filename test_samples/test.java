package com.example.demo;

import java.sql.*;

public class UserService {

    // SQL Injection vulnerability
    public User getUserById(String userId) {
        String query = "SELECT * FROM users WHERE id = '" + userId + "'";
        // This is vulnerable to SQL injection
        return executeQuery(query);
    }

    // Hardcoded credentials
    private static final String DB_PASSWORD = "admin123";

    // Insecure random number generation
    public String generateToken() {
        java.util.Random random = new java.util.Random();
        return String.valueOf(random.nextInt());
    }

    // Command injection vulnerability
    public void executeCommand(String userInput) throws Exception {
        Runtime.getRuntime().exec("ls " + userInput);
    }

    private User executeQuery(String query) {
        // Implementation
        return null;
    }
}
