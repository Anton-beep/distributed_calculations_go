import '../App.css'
import {useState} from "react";
import Auth from "../pkg/Auth";

export const InputExpression = () => {
    const [expression, setExpression] = useState('')
    const [message, setMessage] = useState('')
    const [error, setError] = useState(false)

    const handleChange = (event) => {
        const regex = /^[0-9+\-*/(). ]*$/;
        if (regex.test(event.target.value)) {
            setExpression(event.target.value);
            setMessage('');
            setError(false);
        } else {
            setError(true);
            setMessage("Invalid input. Only numbers and +, /, -, *, ), (, . are allowed.");
        }
    }

    const showMessage = () => {
        if (message === '') {
            return null;
        } else {
            if (error) {
                return <div className="alert alert-danger" role="alert"> {message} </div>
            }
            return <div className="alert alert-success" role="alert"> {message} </div>
        }
    }

    const handleSubmit = (event) => {
        event.preventDefault();
        if (message === '') {
            Auth.axiosInstance.post("/expression", {expression: expression})
                .then(response => {
                    if (response.status === 200) {
                        setError(false);
                        setMessage('Expression added successfully');
                    } else {
                        setError(true);
                        setMessage('Error adding expression');
                    }
                })
                .catch(err => {
                    setError(true);
                    setMessage('Error adding expression');
                });
        }
    }

    return (
        <>
            <h1>Input Expression</h1>
            <p>Supported operations: +, -, *, /</p>
            <form onSubmit={handleSubmit}>
                <label>
                    <input type="text" name="expression" onChange={handleChange} className="form-control" placeholder="(2 + 2) + (2 + 2)"/>
                </label>
                <input type="submit" value="Submit" className="btn btn-secondary"/>
                {showMessage()}
            </form>
        </>
    )
}