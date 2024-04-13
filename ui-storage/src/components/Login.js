import React, {useState} from 'react';
import Auth from '../pkg/Auth';
import Cookies from "js-cookie";

export const Login = () => {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const [error, setError] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        Auth.axiosInstance.post('/login', {
            "login": login,
            "password": password
        })
            .then(response => {
                console.log(response.data.access)
                if (response.status === 200) {
                    Cookies.set('token', response.data.access);
                    setMessage("Success");
                    setError(false)
                    window.location = '/'
                }
            })
            .catch(
                (error) => {
                    setMessage("Invalid login or password");
                    setError(true)
                }
            );
    };

    return (
        <div>
            <h1 className="container">Login</h1>
            <form onSubmit={handleSubmit} className="container">
                <label>
                    Login:
                    <input type="text" value={login} onChange={(e) => setLogin(e.target.value)} required/>
                </label>
                <label>
                    Password:
                    <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required/>
                </label>
                <input type="submit" value="Submit"/>
            </form>
            {message === "" ? null : <div className={error ? "alert alert-danger" : "alert alert-success"}>
                {message}
            </div>}
        </div>
    );
};