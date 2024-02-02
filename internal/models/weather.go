package models

type WeatherApiResponse struct {
    Location struct {
        Name string `json:"name"`
    } `json:"location"`
    Current struct {
        TempC float64 `json:"temp_c"`
        Condition struct {
            Text string `json:"text"`
        } `json:"condition"`
    } `json:"current"`
}