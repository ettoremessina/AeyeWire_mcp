import React, { useState } from 'react';

// XSS vulnerability - dangerouslySetInnerHTML without sanitization
const UserProfile: React.FC = () => {
    const [userBio, setUserBio] = useState('');

    // Vulnerable to XSS
    return (
        <div>
            <div dangerouslySetInnerHTML={{ __html: userBio }} />
        </div>
    );
};

// Hardcoded API key
const API_KEY = 'sk-1234567890abcdef';

// Insecure local storage of sensitive data
const saveAuthToken = (token: string) => {
    localStorage.setItem('authToken', token);
};

// Missing input validation
const LoginForm: React.FC = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const handleSubmit = async () => {
        // No validation, directly sending to API
        await fetch('https://api.example.com/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });
    };

    return (
        <form onSubmit={handleSubmit}>
            <input value={username} onChange={e => setUsername(e.target.value)} />
            <input value={password} onChange={e => setPassword(e.target.value)} />
            <button type="submit">Login</button>
        </form>
    );
};

export { UserProfile, LoginForm };
