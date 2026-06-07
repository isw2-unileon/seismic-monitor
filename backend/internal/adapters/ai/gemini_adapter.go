package ai

import (
	"context"
	"fmt"
	"log"
	"seismic-monitor/backend/internal/models"

	"google.golang.org/genai"
)

type GeminiAdapter struct {
	APIKey string
}

func (a *GeminiAdapter) GenerateSafetyAdvice(ctx context.Context, sismo models.Feature) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  a.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Printf("[Gemini Adapter] Error inicializando cliente: %v", err)
		return "", err
	}

	depth := 0.0
	if len(sismo.Geometry.Coordinates) >= 3 {
		depth = sismo.Geometry.Coordinates[2]
	}

	prompt := fmt.Sprintf(
		"Actúa como un experto en gestión de catástrofes. Se ha detectado un sismo en %s de magnitud %.1f y profundidad %.1f km. "+
			"Proporciona un análisis de riesgo muy breve (máximo 2 líneas) según la situación y zona donde se produce y 3 consejos de seguridad específicos para esta situación.",
		sismo.Info.Place, sismo.Info.Mag, depth,
	)

	resp, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		log.Printf("[⚠️ Gemini API Error] Fallo al generar contenido: %v", err)
		return "Mantente en un lugar seguro y sigue las instrucciones de las autoridades locales.", nil
	}

	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		return resp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "Sin consejos adicionales disponibles en este momento.", nil
}
