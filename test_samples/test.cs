using System;
using System.Data.SqlClient;
using System.Diagnostics;

namespace SecurityDemo
{
    public class UserService
    {
        // SQL Injection vulnerability
        public User GetUserById(string userId)
        {
            string query = "SELECT * FROM Users WHERE Id = '" + userId + "'";
            // Vulnerable to SQL injection
            return ExecuteQuery(query);
        }

        // Hardcoded credentials
        private const string ApiKey = "sk-1234567890abcdef";

        // Command injection vulnerability
        public void ExecuteCommand(string userInput)
        {
            Process.Start("cmd.exe", "/c " + userInput);
        }

        // Weak cryptography
        public string HashPassword(string password)
        {
            using (var md5 = System.Security.Cryptography.MD5.Create())
            {
                // MD5 is weak and deprecated
                byte[] bytes = md5.ComputeHash(System.Text.Encoding.UTF8.GetBytes(password));
                return Convert.ToBase64String(bytes);
            }
        }

        private User ExecuteQuery(string query)
        {
            // Implementation
            return null;
        }
    }
}
