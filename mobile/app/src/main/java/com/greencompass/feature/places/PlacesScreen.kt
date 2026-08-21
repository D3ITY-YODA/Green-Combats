package com.greencompass.feature.places

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Add
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*
import com.greencompass.domain.model.UserPlace
import com.greencompass.navigation.AppRoute

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun PlacesScreen(
    onBack: () -> Unit,
    onAddPlace: () -> Unit,
    onOpenPlace: (String) -> Unit
) {
    val places = remember {
        listOf(
            UserPlace("1", "Lower Valley", true),
            UserPlace("2", "East Ward", false),
            UserPlace("3", "North Basin", false)
        )
    }

    GreenCompassScaffold(
        title = "Your places",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg)
        ) {
            Text(
                text = "Choose the places you want to follow.",
                style = GreenCompassTypography.bodyMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            LazyColumn(
                verticalArrangement = Arrangement.spacedBy(AppSpacing.sm),
                modifier = Modifier.weight(1f)
            ) {
                items(places) { place ->
                    PlaceCard(place = place, onClick = { onOpenPlace(place.id) })
                }
            }

            Spacer(modifier = Modifier.height(AppSpacing.lg))

            SecondaryButton(
                text = "Add a place",
                onClick = onAddPlace,
                modifier = Modifier.padding(bottom = AppSpacing.xxl)
            )
        }
    }
}

@Composable
private fun PlaceCard(place: UserPlace, onClick: () -> Unit) {
    Card(
        modifier = Modifier.fillMaxWidth().clickable { onClick() },
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = Color.White),
        border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)
    ) {
        Row(
            modifier = Modifier.padding(AppSpacing.lg).fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column {
                Text(text = place.name, style = GreenCompassTypography.titleMedium, color = GreenCompassColors.Charcoal)
                if (place.isPrimary) {
                    Spacer(modifier = Modifier.height(AppSpacing.xxs))
                    Text(text = "Primary place", style = GreenCompassTypography.labelSmall, color = GreenCompassColors.ForestGreen)
                }
            }
            TextLinkButton(text = "Open", onClick = onClick)
        }
    }
}
