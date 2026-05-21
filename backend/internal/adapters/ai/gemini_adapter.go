package ai

import (
	"context"
	"fmt"
	"seismic-monitor/backend/internal/models"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiAdapter struct {
	APIKey string
}

func (a *GeminiAdapter) GenerateSafetyAdvice(ctx context.Context, sismo models.Feature) (string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(a.APIKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	depth := 0.0
	if len(sismo.Geometry.Coordinates) >= 3 {
		depth = sismo.Geometry.Coordinates[2]
	}

	prompt := fmt.Sprintf(
		"Actúa como un experto en gestión de catástrofes. Se ha detectado un sismo en %s de magnitud %.1f y profundidad %.1f km. "+
			"Proporciona un análisis de riesgo muy breve (máximo 2 líneas) según la situación y zona donde se produce y 3 consejos de seguridad específicos para esta situación.",
		sismo.Info.Place, sismo.Info.Mag, depth,
	)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "Mantente en un lugar seguro y sigue las instrucciones de las autoridades locales.", nil
	}

	if resp != nil && len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil && len(resp.Candidates[0].Content.Parts) > 0 {
		if part, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(part), nil
		}
	}

	return "Sin consejos adicionales disponibles en este momento.", nil
}
