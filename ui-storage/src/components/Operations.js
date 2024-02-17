import '../App.css'
import {useEffect, useState} from "react";

export const Operations = () => {
    const [operations, setOperations] = useState(null)
    const [messages, setMessages] = useState(null)
    const [mainMessage, setMainMessage] = useState('')
    const [mainError, setMainError] = useState(false)

    useEffect(() => {
        let addr;
        if (process.env.REACT_APP_STORAGE_API_URL === undefined) {
            addr = process.env.REACT_APP_STORAGE_API_URL + "/getOperationsAndTimes"
        } else {
            addr = "http://localhost:8080/api/v1/getOperationsAndTimes"
        }
        fetch(addr)
            .then(response => response.json())
            .then(data => {
                setOperations(data.data)
            })
            .catch(err => {
                setMainError(true);
                setMainMessage('Error getting operations and times')
            });
    }, []);

    const showOperations = () => {
        if (operations === null) {
            return null;
        }

        const table = []
        for (const operation in operations) {
            if (operations.hasOwnProperty(operation)) {
                table.push(
                    <ul className="list-group list-group-horizontal" key={operation}>
                        <li className="list-group-item list-group-item-primary">{operation}</li>
                        <li className="list-group-item list-group-item-primary">
                            <input className="form-control" defaultValue={operations[operation]} onChange={handleChange}
                                   itemID={operation}/>
                            {showMessage(operation)}
                        </li>
                    </ul>
                )
            }
        }

        return table;
    }

    const showMessage = (operation) => {
        if (messages === null) {
            return null;
        }
        if (messages[operation] === '') {
            return null;
        }
        return <div className="alert alert-danger" role="alert"> {messages[operation]} </div>
    }

    const handleChange = (event) => {
        // check value
        const regex = /^[0-9]*$/;
        if (!regex.test(event.target.value)) {
            const newMessages = {...messages}
            newMessages[event.target.getAttribute("itemID")] = "Invalid input. Only numbers are allowed.";
            setMessages(newMessages)
            return;
        }

        const newOperations = {...operations}
        newOperations[event.target.getAttribute("itemID")] = parseInt(event.target.value);
        setOperations(newOperations)
    }

    const handleSubmit = (event) => {
        event.preventDefault();
        let addr;
        if (process.env.REACT_APP_STORAGE_API_URL === undefined) {
            addr = process.env.REACT_APP_STORAGE_API_URL + "/postOperationsAndTimes"
        } else {
            addr = "http://localhost:8080/api/v1/postOperationsAndTimes"
        }
        fetch(addr, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(operations)
        })
            .then(response => {
                if (response.status === 200) {
                    setMainError(false);
                    setMainMessage('Operations and times saved successfully');
                } else {
                    setMainError(true);
                    setMainMessage('Error saving operations and times');
                }
            })
            .catch(err => {
                setMainError(true);
                setMainMessage('Error saving operations and times')
            });
    }

    const showMainMessage = () => {
        if (mainMessage === '') {
            return null;
        }
        if (mainError) {
            return <div className="alert alert-danger" role="alert"> {mainMessage} </div>
        }
        return <div className="alert alert-success" role="alert"> {mainMessage} </div>
    }

    return (
        <>
            <h1>View Operations And Execution Times</h1>
            <ul className="list-group list-group-horizontal">
                <li className="list-group-item">Operation</li>
                <li className="list-group-item">Time (in milliseconds)</li>
            </ul>
            {showOperations()}
            <button type="submit" className="btn btn-secondary" onClick={handleSubmit}>Save</button>
            {showMainMessage()}
        </>
    )
}