import axios from "axios";

export const sendOrderToRabit = (orderUuid) => {
    const formData = new FormData();
    formData.append('uuid', orderUuid);
    return axios.post('http://localhost:8080/api/rabbitmq/start', formData)
}
