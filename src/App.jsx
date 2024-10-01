import { useEffect } from "react"

function App() {
 useEffect(()=>{
  const fetchData = async () => {
    try {
      const response = await fetch('/api/test');
      if (!response.ok) {
        throw new Error('Network response was not ok');
      }
      console.log(await response.text()); 
    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);
    }
  };
  fetchData();
 })
 return <div>
  <h1>
    Welcome to the App!
  </h1>
 </div> 
}

export default App
