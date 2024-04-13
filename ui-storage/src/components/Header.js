import { Link } from 'react-router-dom';
import '../App.css'
import {useEffect, useState} from "react";
import Auth from "../pkg/Auth";

export const Header = () => {
    const [user, setUser] = useState(null);

    // get user information
    useEffect(() => {
        Auth.axiosInstance.get('/getUser')
            .then(response => {
                setUser(response.data.login);
            })
            .catch(err => console.log(err));
    }, [user]);

    let content;
    if (user !== null) {
        content = (
            <header className="header">
                <Link to="/" className="btn btn-primary">Home</Link>
                <Link to="/inputExpression" className="btn btn-primary">Input New Expression</Link>
                <Link to="/viewExpressions" className="btn btn-primary">View All Expressions</Link>
                <Link to="/operations" className="btn btn-primary">View Operations And Execution Times</Link>
                <Link to="/computingPowers" className="btn btn-primary">View Computing Powers</Link>
                <Link to="/logout" className="btn btn-primary">Logout</Link>
            </header>
        )
    } else {
        content = (
            <header className="header">
                <Link to="/login" className="btn btn-primary">Login</Link>
                <Link to="/register" className="btn btn-primary">Register</Link>
            </header>
        )
    }

    return content
}