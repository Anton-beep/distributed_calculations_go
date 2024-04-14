import {useEffect, useState} from "react";
import Auth from "../pkg/Auth";
import Cookies from "js-cookie";

export function Profile() {
    const [oldPassword, setOldPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [login, setLogin] = useState('');
    const [message, setMessage] = useState('');
    const [error, setError] = useState(false);

    useEffect(() => {
        Auth.axiosInstance.get('/getUser')
            .then(response => {
                setLogin(response.data.login);
            })
            .catch(err => console.log(err));
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (oldPassword === '') {
            setMessage("Enter current password");
            setError(true);
            return;
        }
        let sendObject = {
            "old_password": oldPassword,
        }
        if (confirmPassword !== '' || newPassword !== '') {
            if (newPassword !== confirmPassword) {
                setMessage("New Passwords do not match");
                setError(true);
                return;
            }
            sendObject.password = newPassword;
        }
        if (login !== '') {
            sendObject.login = login;
        }

        if (sendObject.password === undefined && sendObject.login === undefined) {
            setMessage("No changes");
            setError(true);
            return;
        }
        Auth.axiosInstance.post('/updateUser', sendObject)
            .then(response => {
                if (response.status === 200) {
                    console.log(response.data)
                    Cookies.set('token', response.data.access);
                    setMessage("Success");
                    setError(false);
                    window.location = '/profile'
                }
            })
            .catch(
                (error) => {
                    console.log(error);
                    setMessage("Invalid data");
                    setError(true);
                }
            );
    }

    return (
        <div>
            <h1>Profile</h1>
            <form onSubmit={handleSubmit}>
                <div>
                    Write your current password if you want to change login or password:
                </div>
                <div className="mb-3">
                    <label htmlFor="oldPassword" className="form-label">Current Password</label>
                    <input type="password" className="form-control" id="oldPassword" value={oldPassword}
                           onChange={e => setOldPassword(e.target.value)}/>
                </div>
                <br></br>
                <div className="mb-3">
                    <label htmlFor="login" className="form-label">Login</label>
                    <input type="text" className="form-control" id="login" value={login}
                           onChange={e => setLogin(e.target.value)}/>
                </div>
                <div className="mb-3">
                    <label htmlFor="newPassword" className="form-label">New Password</label>
                    <input type="password" className="form-control" id="newPassword" value={newPassword}
                           onChange={e => setNewPassword(e.target.value)}/>
                </div>
                <div className="mb-3">
                    <label htmlFor="confirmPassword" className="form-label">Confirm New Password</label>
                    <input type="password" className="form-control" id="confirmPassword" value={confirmPassword}
                           onChange={e => setConfirmPassword(e.target.value)}/>
                </div>
                <button type="submit" className="btn btn-primary">Update</button>
            </form>
            {message === "" ? null : <div className={error ? "alert alert-danger" : "alert alert-success"}>
                {message}
            </div>}
        </div>
    )
}