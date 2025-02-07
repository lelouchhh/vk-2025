import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { Container, Typography, TextField, Button, Link, Box, Alert } from '@mui/material';

const Login: React.FC = () => {
    const [loginValue, setLogin] = useState('');
    const [passwordValue, setPassword] = useState('');
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            const response = await axios.post('http://localhost:8080/login', { login: loginValue, password: passwordValue });
            localStorage.setItem('token', response.data.token);
            setError(null);
            navigate('/containers');
        } catch (err) {
            setError('Invalid login or password');
        }
    };

    return (
        <Container maxWidth="sm">
            <Box sx={{ mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                <Typography variant="h4" gutterBottom>
                    Login
                </Typography>
                {error && <Alert severity="error">{error}</Alert>}
                <Box component="form" onSubmit={handleLogin} sx={{ width: '100%', mt: 3 }}>
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        label="Login"
                        value={loginValue}
                        onChange={(e) => setLogin(e.target.value)}
                    />
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        label="Password"
                        type="password"
                        value={passwordValue}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <Button type="submit" fullWidth variant="contained" sx={{ mt: 3, mb: 2 }}>
                        Login
                    </Button>
                    <Link href="/register" variant="body2">
                        Don't have an account? Register here
                    </Link>
                </Box>
            </Box>
        </Container>
    );
};

export default Login;