# go-idw

**go-idw** is a lightweight and reusable Go library for spatial interpolation using the [Inverse Distance Weighting (IDW)](https://en.wikipedia.org/wiki/Inverse_distance_weighting) method. It allows you to estimate environmental data such as temperature, humidity, rainfall, and wind speed at any geographical coordinate based on known data points.

## 🌐 What is IDW?

Inverse Distance Weighting is a deterministic method for spatial interpolation that estimates the value of an unknown point using the values from surrounding known points. The influence of each known point is inversely proportional to its distance from the target point.

## 🚀 Features

- Estimate **4 fixed environmental variables**:
  - Noon wind speed (km/h)
  - Noon temperature (°C)
  - Noon relative humidity (%)
  - 24-hour rainfall (mm)
- Uses the **Haversine formula** (https://en.wikipedia.org/wiki/Haversine_formula) to calculate distance over the Earth's surface
- Easily embeddable and customizable
- Designed for **wildfire risk modeling**, **climate analysis**, and **map visualization**

## 📦 Installation

```bash
go get github.com/cosholyknight/go-idw
```
