import {useEffect, useState} from "react";

function Hello(){
    function byeFn() {
        console.log("bye :(")
    }
    function hiFn(){
        console.log("created :)");
        return byeFn;
    }
    useEffect(hiFn, []);
    return <h1>Hello</h1>;
}

function App() {
    const [showing, setShwoing] = useState(false);
    const onClick = () => setShwoing((prev) => !prev)

    return (
        <div>
            {showing ? <Hello /> : null}
            <button onClick={onClick}>{showing ? "Hide" : "Show"}</button>
        </div>
    )
}

export default App;
