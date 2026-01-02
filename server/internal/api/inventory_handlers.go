// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// InventoryResponse represents the inventory API response
type InventoryResponse struct {
	Projects []ProjectInventoryResponse `json:"projects"`
}

// ProjectInventoryResponse represents project inventory in API response
type ProjectInventoryResponse struct {
	Project ProjectResponse           `json:"project"`
	Apps    []AppInventoryResponse    `json:"apps"`
}

// AppInventoryResponse represents app inventory in API response
type AppInventoryResponse struct {
	App      AppResponse          `json:"app"`
	Versions []VersionResponse    `json:"versions"`
}

// handleGetInventory godoc
// @Summary      Get complete inventory
// @Description  Get complete inventory data including all projects, apps, and versions
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  InventoryResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /admin/inventory [get]
func (h *Handler) handleGetInventory(c *gin.Context) {
	inventory, err := h.inventoryService.GetCompleteInventory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to API response format
	response := InventoryResponse{
		Projects: make([]ProjectInventoryResponse, len(inventory.Projects)),
	}

	for i, projectInventory := range inventory.Projects {
		appInventories := make([]AppInventoryResponse, len(projectInventory.Apps))
		for j, appInventory := range projectInventory.Apps {
			versions := make([]VersionResponse, len(appInventory.Versions))
			for k, version := range appInventory.Versions {
				versions[k] = VersionResponse{
					ID:          version.ID,
					AppID:       version.AppID,
					Version:     version.Hash,
					IsPublished: version.IsPublished,
					CreatedAt:   version.CreatedAt.Format(time.RFC3339),
				}
			}

			appInventories[j] = AppInventoryResponse{
				App: AppResponse{
					ID:        appInventory.App.ID,
					ProjectID: appInventory.App.ProjectID,
					Name:      appInventory.App.Name,
					CreatedAt: appInventory.App.CreatedAt.Format(time.RFC3339),
				},
				Versions: versions,
			}
		}

		response.Projects[i] = ProjectInventoryResponse{
		Project: ProjectResponse{
			ID:        projectInventory.Project.ID,
			Name:      projectInventory.Project.Name,
			CreatedAt: projectInventory.Project.CreatedAt.Format(time.RFC3339),
		},
			Apps: appInventories,
		}
	}

	c.JSON(http.StatusOK, response)
}

// handleGetProjectInventory godoc
// @Summary      Get project inventory
// @Description  Get inventory data for a specific project
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        project  path      string  true  "Project name"
// @Success      200      {object}  ProjectInventoryResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      404      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Security     Bearer
// @Router       /admin/inventory/{project} [get]
func (h *Handler) handleGetProjectInventory(c *gin.Context) {
	projectName := c.Param("project")

	projectInventory, err := h.inventoryService.GetProjectInventory(projectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to API response format
	appInventories := make([]AppInventoryResponse, len(projectInventory.Apps))
	for j, appInventory := range projectInventory.Apps {
		versions := make([]VersionResponse, len(appInventory.Versions))
		for k, version := range appInventory.Versions {
			versions[k] = VersionResponse{
				ID:          version.ID,
				AppID:       version.AppID,
				Version:     version.Hash,
				IsPublished: version.IsPublished,
				CreatedAt:   version.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			}
		}

		appInventories[j] = AppInventoryResponse{
			App: AppResponse{
				ID:        appInventory.App.ID,
				ProjectID: appInventory.App.ProjectID,
				Name:      appInventory.App.Name,
				CreatedAt: appInventory.App.CreatedAt.Format(time.RFC3339),
			},
			Versions: versions,
		}
	}

	response := ProjectInventoryResponse{
		Project: ProjectResponse{
			ID:        projectInventory.Project.ID,
			Name:      projectInventory.Project.Name,
			CreatedAt: projectInventory.Project.CreatedAt.Format(time.RFC3339),
		},
		Apps: appInventories,
	}

	c.JSON(http.StatusOK, response)
}

// handleGetInventorySummary godoc
// @Summary      Get inventory summary
// @Description  Get summary statistics for the inventory
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  InventorySummaryResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     Bearer
// @Router       /admin/inventory/summary [get]
func (h *Handler) handleGetInventorySummary(c *gin.Context) {
	summary, err := h.inventoryService.GetInventorySummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

