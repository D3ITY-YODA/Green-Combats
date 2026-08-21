package com.greencompass.feature.today

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.CheckCircle
import androidx.compose.material.icons.outlined.Warning
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.greencompass.core.ui.*
import com.greencompass.domain.model.StatusType
import com.greencompass.domain.model.TodayData

@Composable
fun TodayRoute(viewModel: TodayViewModel = hiltViewModel()) {
    val state by viewModel.uiState.collectAsState()
    TodayScreen(state = state, onRetry = { viewModel.retry() })
}

@Composable
fun TodayScreen(state: TodayUiState, onRetry: () -> Unit) {
    GreenCompassScaffold(
        title = "Today",
        showPlaceSwitcher = true,
        placeName = state.data?.placeName ?: "Lower Valley"
    ) { paddingValues ->
        Column(modifier = Modifier.fillMaxSize().padding(paddingValues)) {
            if (state.isOffline) {
                OfflineBanner()
            }

            if (state.isLoading) {
                Column(modifier = Modifier.fillMaxSize().padding(AppSpacing.lg)) {
                    LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(120.dp))
                    Spacer(Modifier.height(AppSpacing.md))
                    LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(80.dp))
                    Spacer(Modifier.height(AppSpacing.md))
                    LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(80.dp))
                }
            } else if (state.data != null) {
                TodayContent(data = state.data, isOffline = state.isOffline, onRetry = onRetry)
            }
        }
    }
}

@Composable
private fun TodayContent(data: TodayData, isOffline: Boolean, onRetry: () -> Unit) {
    LazyColumn(
        modifier = Modifier.fillMaxSize().padding(AppSpacing.lg),
        verticalArrangement = Arrangement.spacedBy(AppSpacing.md)
    ) {
        // Status Card based on state
        item {
            when (data.status.type) {
                StatusType.DELAYED -> DelayedStateCard(status = data.status, onRetry = onRetry)
                else -> TodayStatusCard(status = data.status)
            }
        }

        // Important Updates (if any)
        if (data.updates.isNotEmpty()) {
            items(data.updates) { update ->
                ImportantUpdateCard(update = update)
            }
            if (data.updates.size > 1) {
                item {
                    SecondaryButton(
                        text = "View all updates",
                        onClick = { /* Navigate to Updates */ },
                        modifier = Modifier.padding(top = AppSpacing.xs)
                    )
                }
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
        
        // Bottom padding for navigation bar
        item {
            Spacer(modifier = Modifier.height(AppSpacing.xxxl))
        }
    }
}

@Composable
private fun TodayStatusCard(status: com.greencompass.domain.model.TodayStatus) {
    val (icon, titleColor) = when (status.type) {
        StatusType.IMPORTANT -> Icons.Outlined.Warning to GreenCompassColors.ImportantAmber
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
private fun DelayedStateCard(status: com.greencompass.domain.model.TodayStatus, onRetry: () -> Unit) {
    Card(
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Column(modifier = Modifier.padding(AppSpacing.lg), horizontalAlignment = Alignment.CenterHorizontally) {
            Icon(imageVector = Icons.Outlined.Warning, contentDescription = null, tint = GreenCompassColors.DelayedGrey, modifier = Modifier.size(48.dp))
            Spacer(Modifier.height(AppSpacing.md))
            Text(text = "Information delayed", style = GreenCompassTypography.headlineMedium, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
            Spacer(Modifier.height(AppSpacing.sm))
            Text(text = "The latest update for Lower Valley\nis not available yet.", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
            Spacer(Modifier.height(AppSpacing.sm))
            Text(text = "Last reliable update: ${status.updatedAt}", style = GreenCompassTypography.labelMedium, color = GreenCompassColors.DelayedGrey)
            Spacer(Modifier.height(AppSpacing.lg))
            Row(horizontalArrangement = Arrangement.spacedBy(AppSpacing.sm)) {
                SecondaryButton(text = "Try again", onClick = onRetry, modifier = Modifier.weight(1f).height(AppSpacing.huge))
            }
        }
    }
}

@Composable
private fun ImportantUpdateCard(update: com.greencompass.domain.model.PublicUpdate) {
    Card(
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Column(modifier = Modifier.padding(AppSpacing.lg)) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Box(modifier = Modifier.size(8.dp).padding(end = AppSpacing.xs).background(GreenCompassColors.ImportantAmber, shape = CircleShape))
                Text(text = update.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
            }
            Spacer(Modifier.height(AppSpacing.xs))
            Text(text = update.message, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(AppSpacing.sm))
            Text(text = "${update.placeName} · Updated ${update.updatedAt}", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(AppSpacing.sm))
            TextLinkButton(text = "Read update", onClick = { /* Navigate to detail */ })
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
