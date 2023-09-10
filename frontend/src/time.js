function format(time) {
    let month = (time.getMonth() + 1 + "").padStart(2, "0")
    let day = (time.getDate() + "").padStart(2, "0")
    let year = (time.getFullYear() + "").padStart(4, "0");
    let hour = (time.getHours() + "").padStart(2, "0")
    let minute = (time.getMinutes() + "").padStart(2, "0")
    let second = (time.getSeconds() + "").padStart(2, "0")

    let timezone = "CDT"

    return `${month}/${day}/${year}|${hour}:${minute}:${second}|${timezone}`
}

export default { format }