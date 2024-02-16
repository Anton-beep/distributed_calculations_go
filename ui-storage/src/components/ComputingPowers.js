import '../App.css'
import {useEffect, useState} from "react";

export const ComputingPowers = () => {
    const [servers, setServers] = useState(null)

    useEffect(() => {
        fetch(process.env.REACT_APP_STORAGE_API_URL + "/getComputingPowers")
            .then(response => response.json())
            .then(data => {
                setServers(data.servers)
                console.log(data)
            })
            .catch(err => console.log(err));
    }, []);

    return (
        <>
            <title>View Computing Powers</title>
            <h1>View Computing Powers</h1>
            <div className="scrollable-div">
                <ul className="list-group list-group-horizontal">
                    <li className="list-group-item">Server Name</li>
                    <li className="list-group-item">Calculated Expressions IDs</li>
                </ul>
                {servers !== null ? servers.map((server, index) => (
                        <ul className="list-group list-group-horizontal" key={index}>
                            <li className="list-group-item list-group-item-primary">{server.server_name}</li>
                            <li className="list-group-item list-group-item-primary">{server.calculated_expressions.join("; ")}</li>
                        </ul>
                    )
                ) : null}
            </div>
        </>
    )
}