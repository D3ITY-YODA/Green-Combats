package com.greencompass.feature.today

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.CheckCircle
import androidx.compose.material.icons.outlined.Warning
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.greencompass.core.ui.*
import com.greencompass.domain.model.StatusType
import com.greencompass.domain.model.TodayData

@Composable
fun TodayRoute(viewModel: TodayViewModel = hiltViewModel()) {
    val state by viewModel.uiState.collectAsState()
    TodayScreen(state = state)
}

@Composable
fun TodayScreen(state: TodayUiState) {
    GreenCompassScaffold(
        title = "Today",
        showPlaceSwitcher = true,
        placeName = state.data?.placeName ?: "Lower Valley"
    ) { paddingValues ->
        if (state.isLoading) {
            Column(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(AppSpacing.lg)) {
                LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(120.dp))
                Spacer(Modifier.height(AppSpacing.md))
                LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(80.dp))
            }
        } else if (state.data != null) {
            TodayContent(data = state.data, modifier = Modifier.padding(paddingValues))
        }
    }
}

@Composable
private fun TodayContent(data: TodayData, modifier: Modifier = Modifier) {
    LazyColumn(
        modifier = modifier.fillMaxSize().padding(AppSpacing.lg),
        verticalArrangement = Arrangement.spacedBy(AppSpacing.md)
    ) {
        // Status Card
        item {
            StatusCard(status = data.status)
        }

        // Updates
        if (data.updates.isNotEmpty()) {
            items(data.updates) { update ->
                UpdateCard(update = update)
            }
        }

        // Explore Sections
        if (data.exploreSections.isNotEmpty()) {
            item {
                Text(
                    text = "Explore",
                    style = GreenCompassTypography.titleLarge,
                    color = GreenCompassColors.Charcoal,
                    modifier = Modifier.padding(top = AppSpacing.md, bottom = AppSpacing.sm)
                )
            }
            items(data.exploreSections) { section ->
                ExploreCard(section = section)
            }
        }
    }
}

@Composable
private fun StatusCard(status: com.greencompass.domain.model.TodayStatus) {
    val (icon, titleColor) = when (status.type) {
        StatusType.IMPORTANT -> Icons.Outlined.Warning to GreenCompassColors.ImportantAmber
        StatusType.DELAYED -> Icons.Outlined.Warning to GreenCompassColors.DelayedGrey
        else -> Icons.Outlined.CheckCircle to GreenCompassColors.NormalGreen
    }

    Card(
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Row(modifier = Modifier.padding(AppSpacing.lg), verticalAlignment = Alignment.CenterVertically) {
            Icon(imageVector = icon, contentDescription = null, tint = titleColor, modifier = Modifier.size(32.dp))
            Spacer(Modifier.width(AppSpacing.md))
            Column {
                Text(text = status.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                Text(text = status.message, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
                Text(text = "Updated ${status.updatedAt}", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText, modifier = Modifier.padding(top = AppSpacing.xs))
            }
        }
    }
}

@Composable
private fun UpdateCard(update: com.greencompass.domain.model.PublicUpdate) {
    Card(
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Column(modifier = Modifier.padding(AppSpacing.lg)) {
            Text(text = update.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
            Spacer(Modifier.height(AppSpacing.xs))
            Text(text = update.message, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(AppSpacing.sm))
            Text(text = "${update.placeName} · Updated ${update.updatedAt}", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText)
        }
    }
}

@Composable
private fun ExploreCard(section: com.greencompass.domain.model.ExploreSection) {
    Card(
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Column(modifier = Modifier.padding(AppSpacing.lg)) {
            Text(text = section.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
            Spacer(Modifier.height(AppSpacing.xs))
            Text(text = section.description, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
        }
    }
}
