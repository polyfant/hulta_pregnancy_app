interface WeatherData {
	temperature: number;
	humidity: number;
	precipitation: number;
	season: 'spring' | 'summer' | 'fall' | 'winter';
	airQuality: {
		pm25: number;
		pm10: number;
		no2: number;
	};
	pressure: number; // barometric pressure
}

interface EnvironmentalImpact {
	growthImpact: number; // -1 to 1 scale
	recommendations: string[];
	riskFactors: string[];
}

export const environmentalFactors = {
	async analyzeEnvironmentalImpact(location: string): Promise<WeatherData> {
		// Fetch local weather data (you'll need to add your weather API key)
		const weather = await fetch(
			`https://api.weatherapi.com/v1/forecast.json?key=${process.env.WEATHER_API_KEY}&q=${location}&days=7`
		).then((res) => res.json());

		return {
			temperature: weather.current.temp_c,
			humidity: weather.current.humidity,
			precipitation: weather.forecast.forecastday[0].day.totalprecip_mm,
			season: this.getSeason(new Date()),
			airQuality: {
				pm25: 0,
				pm10: 0,
				no2: 0,
			},
			pressure: 0,
		};
	},

	calculateEnvironmentalImpact(data: WeatherData): EnvironmentalImpact {
		const impact = {
			growthImpact: 0,
			recommendations: [],
			riskFactors: [],
		};

		// Temperature impact
		if (data.temperature < 5) {
			impact.growthImpact -= 0.3;
			impact.recommendations.push(
				'Consider additional shelter or blankets'
			);
			impact.riskFactors.push('Cold stress');
		} else if (data.temperature > 30) {
			impact.growthImpact -= 0.2;
			impact.recommendations.push('Ensure adequate shade and cooling');
			impact.riskFactors.push('Heat stress');
		}

		// Humidity impact
		if (data.humidity > 80) {
			impact.growthImpact -= 0.1;
			impact.recommendations.push('Monitor for respiratory issues');
		}

		return impact;
	},

	async getAirQuality(location: string) {
		const response = await fetch(
			`https://api.openweathermap.org/data/2.5/air_pollution?lat=${location.lat}&lon=${location.lon}&appid=${process.env.WEATHER_API_KEY}`
		);
		return response.json();
	},

	calculateAirQualityImpact(aq: WeatherData['airQuality']) {
		let impact = 0;
		const risks = [];

		if (aq.pm25 > 35) {
			impact -= 0.2;
			risks.push('High PM2.5 levels - Consider indoor time');
		}
		if (aq.no2 > 100) {
			impact -= 0.15;
			risks.push('Elevated NO2 - Monitor respiratory health');
		}

		return { impact, risks };
	},
};
