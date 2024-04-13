import React, {useState} from 'react';
import Auth from '../pkg/Auth';
import Cookies from "js-cookie";

export const Register = () => {
    const [login, setLogin] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const [error, setError] = useState(false);

    const handleSubmit = async (e) => {
        e.preventDefault();
        Auth.axiosInstance.post('/register', {
            "login": login,
            "password": password
        })
            .then(response => {
                console.log(response)
                if (response.status === 200) {
                    Cookies.set('token', response.data.access);
                    setMessage("Success");
                    setError(false)
                    window.location = '/'
                }
            })
            .catch(
                (error) => {
                    console.log(error)
                    if (error.response.status === 409) {
                        setMessage("User with such login already exists");
                        setError(true)
                    }
                }
            );
    };

    return (
        <div>
            <h1 className="container">Register</h1>
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