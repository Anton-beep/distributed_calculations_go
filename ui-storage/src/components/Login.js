import React, { useState } from 'react';
import Auth  from '../pkg/Auth';

export const Login = () => {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        Auth.login(login, password)
            .then(response => {
                console.log(response)
                if (response.status === 200) {
                    setMessage("Success");
                } else {
                    setMessage("Error");
                }
            });
    };

    return (
        <div>
            <h1>Login</h1>
            <form onSubmit={handleSubmit}>
                <label>
                    Username:
                    <input type="text" value={login} onChange={(e) => setLogin(e.target.value)} />
                </label>
                <label>
                    Password:
                    <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
                </label>
                <input type="submit" value="Submit" />
            </form>
        </div>
    );
};