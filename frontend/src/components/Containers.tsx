import React, { useEffect, useState } from 'react';
import axios from 'axios';
import {
    Container,
    Typography,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    Alert,
} from '@mui/material';

const Containers: React.FC = () => {
    const [containers, setContainers] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchContainers = async () => {
            try {
                const token = localStorage.getItem('token');
                const response = await axios.get('http://localhost:8080/protected/containers', {
                    headers: { Authorization: `Bearer ${token}` },
                });
                setContainers(response.data);
            } catch (err) {
                setError('Failed to fetch containers');
            } finally {
                setLoading(false);
            }
        };

        fetchContainers();
    }, []);

    if (loading) return <Typography>Loading...</Typography>;
    if (error) return <Alert severity="error">{error}</Alert>;

    return (
        <Container maxWidth="lg">
            <Typography variant="h4" gutterBottom>
                Containers
            </Typography>

            <TableContainer component={Paper}>
                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>IP Address</TableCell>
                            <TableCell>Last Ping</TableCell>
                            <TableCell>Status</TableCell>
                            <TableCell>Ping Time</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {containers.map((container) => (
                            <TableRow key={container.id}>
                                <TableCell>{container.ip_address}</TableCell>
                                <TableCell>{container.last_ping}</TableCell>
                                <TableCell>{container.status ? 'Online' : 'Offline'}</TableCell>
                                <TableCell>{container.ping_time || 'N/A'}</TableCell>
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </Container>
    );
};

export default Containers;