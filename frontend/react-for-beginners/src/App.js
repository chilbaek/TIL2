import {useEffect, useState} from "react";

function App() {
    const [counter, setValue] = useState(0);
    const [keyword, setKeyword] = useState("")
    const onClick = () => setValue((prev) => prev + 1);
    const onChange = (event) => setKeyword(event.target.value);
    useEffect(() => {
        console.log("I run all the time");
    }, []);
    useEffect(() => {
        // if(keyword !== "" && keyword.length > 5) {
            console.log("I run when 'keyword' changes.", keyword);
        // }
    }, [keyword]);
   useEffect(() => {
       console.log("I run when 'counter' changes.", counter);
    }, [counter]);
  useEffect(() => {
       console.log("I run when 'counter' & 'keyword' changes.", counter, keyword);
    }, [counter, keyword]);
  return (
    <div>
        <input
            value={keyword}
            onChange={onChange}
            type="text"
            placeholder="search here ..."
        />
      <h1>{counter}</h1>
      <button onClick={onClick}>click me </button>
    </div>
  );
}

export default App;
