import { get } from "svelte/store";
import stores from "./stores";

const BASE = import.meta.env.PROD ? "/api" : "http://localhost:3000/api";

async function httpRequest(route, params, method) {
    params.token = get(stores.user).token
    let status = 200;
    return await fetch(
        BASE + route + "?" + new URLSearchParams(params), { method: method })
        .then((response) => {
            status = response.status;
            if (response.headers.get("Content-Length") == 0) {
                return Promise.resolve({})
            }
            if (status == 200) {
                return response.json();
            } else {
                return response.text();
            }
        })
        .then((body) => {
            return {
                status: status,
                body: body,
            };
        });
}

async function GET(route, params) {
    return await httpRequest(route, params, "GET")
}

async function POST(route, params) {
    return await httpRequest(route, params, "POST")
}

async function PUT(route, params) {
    return await httpRequest(route, params, "PUT")
}

async function DELETE(route, params) {
    return await httpRequest(route, params, "DELETE")
}

export default { GET, POST, PUT, DELETE, BASE }