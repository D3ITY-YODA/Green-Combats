package com.greencompass.feature.updates

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.Info
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.greencompass.core.ui.*
import com.greencompass.domain.model.TodayData
import com.greencompass.feature.today.TodayViewModel
import com.greencompass.navigation.AppRoute

@Composable
fun UpdatesScreen(
    viewModel: TodayViewModel = hiltViewModel(),
    onNavigate: (AppRoute) -> Unit
) {
    val state by viewModel.uiState.collectAsState()
    
    GreenCompassScaffold(
        title = "Updates",
        showPlaceSwitcher = true,
        placeName = state.data?.placeName ?: "Lower Valley"
    ) { paddingValues ->
        if (state.isLoading) {
            Column(modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(start = AppSpacing.lg, end = AppSpacing.lg, top = AppSpacing.sm, bottom = AppSpacing.lg)) {
                repeat(3) { LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(100.dp).padding(bottom = AppSpacing.md)) }
            }
        } else {
            val currentData = state.data
            if (currentData != null) {
                UpdatesContent(data = currentData, modifier = Modifier
                    .padding(paddingValues)
                    .padding(start = AppSpacing.lg, end = AppSpacing.lg, top = AppSpacing.sm, bottom = AppSpacing.lg), onNavigate = onNavigate)
            }
        }
    }
}

@Composable
private fun UpdatesContent(data: TodayData, modifier: Modifier = Modifier, onNavigate: (AppRoute) -> Unit) {
    LazyColumn(
        modifier = modifier.fillMaxSize(),
        verticalArrangement = Arrangement.spacedBy(AppSpacing.md)
    ) {
        if (data.updates.isEmpty()) {
            item {
                EmptyStateView(
                    icon = Icons.Outlined.Info,
                    title = "No important updates",
                    message = "There are no important updates for your selected places."
                )
            }
        } else {
            val important = data.updates.filter { it.isImportant }
            val other = data.updates.filter { !it.isImportant }

            if (important.isNotEmpty()) {
                item {
                    Text(text = "Important updates", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(bottom = AppSpacing.sm))
                }
                items(important) { update -> 
                    UpdateCard(update = update, onClick = { onNavigate(AppRoute.UpdateDetail(update.id)) }) 
                }
            }

            if (other.isNotEmpty()) {
                item {
                    Text(text = "Other updates", style = GreenCompassTypography.titleLarge, color = GreenCompassColors.Charcoal, modifier = Modifier.padding(top = AppSpacing.md, bottom = AppSpacing.sm))
                }
                items(other) { update -> 
                    UpdateCard(update = update, onClick = { onNavigate(AppRoute.UpdateDetail(update.id)) }) 
                }
            }
        }
    }
}

@Composable
private fun UpdateCard(update: com.greencompass.domain.model.PublicUpdate, onClick: () -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth().clickable { onClick() },
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Column(modifier = Modifier.padding(AppSpacing.lg).fillMaxWidth()) {
            Text(text = update.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
            Spacer(Modifier.height(AppSpacing.xs))
            Text(text = update.message, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(AppSpacing.sm))
            Text(text = "${update.placeName} · Updated ${update.updatedAt}", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText)
        }
    }
}
