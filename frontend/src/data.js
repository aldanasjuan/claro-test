import {writable} from 'svelte/store'
export const url = "http://localhost:8855"
export const session = writable(null)
export const user = writable(null)


if(window){
    console.log("window")
    let local = window.localStorage.getItem("session")
    if(local){
        session.set(local)
    }
    session.subscribe(v => {
        if(v){
            try{
                let s = v.split(".")[0]
                user.set(JSON.parse(atob(s)))
            }catch(error){
                console.log(error)
            }
            window.localStorage.setItem("session", v)
        }else{
            user.set(null)
            window.localStorage.removeItem("session")
        }
    })
}