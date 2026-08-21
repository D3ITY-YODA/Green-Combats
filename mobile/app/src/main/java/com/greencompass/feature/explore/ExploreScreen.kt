package com.greencompass.feature.explore

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.greencompass.core.ui.*
import com.greencompass.domain.model.ExploreSection
import com.greencompass.feature.today.TodayViewModel

@Composable
fun ExploreScreen(viewModel: TodayViewModel = hiltViewModel()) {
    val state by viewModel.uiState.collectAsState()
    
    GreenCompassScaffold(
        title = "Explore",
        showPlaceSwitcher = true,
        placeName = state.data?.placeName ?: "Lower Valley"
    ) { paddingValues ->
        if (state.isLoading) {
            Column(modifier = Modifier.fillMaxSize().padding(paddingValues).padding(AppSpacing.lg)) {
                repeat(4) { LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(100.dp).padding(bottom = AppSpacing.md)) }
            }
        } else if (state.data != null) {
            ExploreContent(data = state.data!!, modifier = Modifier.padding(paddingValues))
        }
    }
}

@Composable
private fun ExploreContent(data: com.greencompass.domain.model.TodayData, modifier: Modifier = Modifier) {
    LazyColumn(
        modifier = modifier.fillMaxSize().padding(AppSpacing.lg),
        verticalArrangement = Arrangement.spacedBy(AppSpacing.md)
    ) {
        items(data.exploreSections) { section ->
            ExploreCard(section = section)
        }
    }
}

@Composable
private fun ExploreCard(section: ExploreSection) {
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
