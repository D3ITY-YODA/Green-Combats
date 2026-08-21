package com.greencompass.core.ui

import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun GreenCompassScaffold(
    title: String,
    showPlaceSwitcher: Boolean = false,
    placeName: String? = null,
    onPlaceClick: (() -> Unit)? = null,
    navigationIcon: @Composable (() -> Unit)? = null,
    actions: @Composable (() -> Unit)? = null,
    content: @Composable (PaddingValues) -> Unit
) {
    Scaffold(
        topBar = {
            GreenCompassTopBar(
                title = title,
                showPlaceSwitcher = showPlaceSwitcher,
                placeName = placeName,
                onPlaceClick = onPlaceClick,
                navigationIcon = navigationIcon,
                actions = actions
            )
        },
        containerColor = GreenCompassColors.WarmWhite
    ) { paddingValues ->
        Box(modifier = Modifier.padding(paddingValues)) {
            content(paddingValues)
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun GreenCompassTopBar(
    title: String,
    showPlaceSwitcher: Boolean = false,
    placeName: String? = null,
    onPlaceClick: (() -> Unit)? = null,
    navigationIcon: @Composable (() -> Unit)? = null,
    actions: @Composable (() -> Unit)? = null
) {
    TopAppBar(
        title = {
            Column {
                Text(text = title, style = GreenCompassTypography.titleMedium)
                if (showPlaceSwitcher && placeName != null) {
                    PlaceSwitcher(placeName = placeName, onClick = onPlaceClick ?: {})
                }
            }
        },
        navigationIcon = { navigationIcon?.invoke() },
        actions = { actions?.invoke() },
        colors = TopAppBarDefaults.topAppBarColors(
            containerColor = GreenCompassColors.WarmWhite,
            titleContentColor = GreenCompassColors.Charcoal,
            navigationIconContentColor = GreenCompassColors.Charcoal
        )
    )
}

@Composable
fun PlaceSwitcher(placeName: String, onClick: () -> Unit) {
    TextButton(onClick = onClick, contentPadding = PaddingValues(0.dp)) {
        Text(
            text = "$placeName ⌄",
            style = GreenCompassTypography.labelMedium,
            color = GreenCompassColors.MutedText
        )
    }
}
