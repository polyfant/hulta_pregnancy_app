
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  TextInput, 
  PasswordInput, 
  Button, 
  Paper, 
  Title, 
  Container 
} from '@mantine/core';

const LoginPage: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });
      
      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('token', data.token);
        navigate('/dashboard');
      }
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  return (
    <Container size={420} my={40}>
      <Title ta="center" order={1}>
        Welcome to Hulta Pregnancy App
      </Title>

      <Paper withBorder shadow="md" p={30} mt={30} radius="md">
        <form onSubmit={handleLogin}>
          <TextInput
            label="Email"
            placeholder="your@email.com"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <PasswordInput
            label="Password"
            placeholder="Your password"
            required
            mt="md"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <Button type="submit" fullWidth mt="xl">
            Sign in
          </Button>
        </form>
      </Paper>
    </Container>
  );
};

export default LoginPage;