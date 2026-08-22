package com.greencompass.feature.today

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.DateRange
import androidx.compose.material.icons.filled.Info
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Place
import androidx.compose.material.icons.filled.Star
import androidx.compose.material.icons.filled.Warning
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.greencompass.core.ui.*
import com.greencompass.domain.model.ExploreSection
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
                Column(modifier = Modifier
                    .fillMaxSize()
                    .padding(start = AppSpacing.lg, end = AppSpacing.lg, top = AppSpacing.sm, bottom = AppSpacing.lg)) {
                    LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(120.dp))
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
        modifier = Modifier
            .fillMaxSize()
            .padding(start = AppSpacing.lg, end = AppSpacing.lg, top = AppSpacing.sm, bottom = AppSpacing.lg),
        verticalArrangement = Arrangement.spacedBy(AppSpacing.md)
    ) {
        item {
            when (data.status.type) {
                StatusType.DELAYED -> DelayedStateCard(status = data.status, onRetry = onRetry)
                else -> TodayStatusCard(status = data.status)
            }
        }

        if (data.updates.isNotEmpty()) {
            items(data.updates) { update ->
                ImportantUpdateCard(update = update)
            }
            if (data.updates.size > 1) {
                item {
                    SecondaryButton(
                        text = "View all updates",
                        onClick = { },
                        modifier = Modifier.padding(top = AppSpacing.xs)
                    )
                }
            }
        }

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
        
        item {
            Spacer(modifier = Modifier.height(AppSpacing.xxxl))
        }
    }
}

@Composable
private fun TodayStatusCard(status: com.greencompass.domain.model.TodayStatus) {
    val (icon, titleColor) = when (status.type) {
        StatusType.IMPORTANT -> Icons.Filled.Warning to GreenCompassColors.ImportantAmber
        else -> Icons.Filled.CheckCircle to GreenCompassColors.NormalGreen
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
            Icon(imageVector = Icons.Filled.Warning, contentDescription = null, tint = GreenCompassColors.DelayedGrey, modifier = Modifier.size(48.dp))
            Spacer(Modifier.height(AppSpacing.md))
            Text(text = "Information delayed", style = GreenCompassTypography.headlineMedium, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center)
            Spacer(Modifier.height(AppSpacing.sm))
            Text(text = "The latest update for Lower Valley\nis not available yet.", style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText, textAlign = TextAlign.Center)
            Spacer(Modifier.height(AppSpacing.lg))
            SecondaryButton(text = "Try again", onClick = onRetry, modifier = Modifier.height(AppSpacing.huge))
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
                Box(modifier = Modifier.size(8.dp).background(GreenCompassColors.ImportantAmber, shape = CircleShape))
                Spacer(Modifier.width(AppSpacing.xs))
                Text(text = update.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
            }
            Spacer(Modifier.height(AppSpacing.xs))
            Text(text = update.message, style = GreenCompassTypography.bodyMedium, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(AppSpacing.sm))
            Text(text = "${update.placeName} · Updated ${update.updatedAt}", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(AppSpacing.sm))
            TextLinkButton(text = "Read update", onClick = { })
        }
    }
}

private fun getIconForSection(key: String) = when (key) {
    "local" -> Icons.Filled.Info
    "seasonal" -> Icons.Filled.DateRange
    "water" -> Icons.Filled.Place
    "land" -> Icons.Filled.Place
    "food" -> Icons.Filled.Star
    "community" -> Icons.Filled.Person
    else -> Icons.Filled.Info
}

@Composable
private fun ExploreCard(section: ExploreSection) {
    Card(
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Row(modifier = Modifier.padding(AppSpacing.lg).fillMaxWidth(), verticalAlignment = Alignment.CenterVertically) {
            Box(
                modifier = Modifier
                    .size(48.dp)
                    .background(GreenCompassColors.SoftSage, shape = CircleShape),
                contentAlignment = Alignment.Center
            ) {
                Icon(
                    imageVector = getIconForSection(section.key),
                    contentDescription = null,
                    tint = GreenCompassColors.ForestGreen,
                    modifier = Modifier.size(24.dp)
                )
            }
            Spacer(Modifier.width(AppSpacing.md))
            Column(modifier = Modifier.weight(1f)) {
                Text(text = section.title, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                Spacer(Modifier.height(AppSpacing.xxs))
                Text(text = section.description, style = GreenCompassTypography.bodySmall, color = GreenCompassColors.MutedText)
            }
        }
    }
}
